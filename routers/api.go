package routers

import (
	"io"
	"net/http"

	"ginskeleton/app/global/consts"
	"ginskeleton/app/global/variable"
	"ginskeleton/app/http/middleware/cors"
	validatorFactory "ginskeleton/app/http/validator/core/factory"
	"ginskeleton/app/utils/gin_release"

	"github.com/natefinch/lumberjack"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// 该路由主要设置门户类网站等前台路由

func InitApiRouter() *gin.Engine {
	var router *gin.Engine
	// 非调试模式（生产模式） 日志写到日志文件
	if variable.ConfigYml.GetBool("AppDebug") == false {
		// 1.gin自行记录接口访问日志，不需要nginx，如果开启以下3行，那么请屏蔽第 34 行代码
		// gin.DisableConsoleColor()
		// f, _ := os.Create(variable.BasePath + variable.ConfigYml.GetString("Logs.GinLogName"))
		// gin.DefaultWriter = io.MultiWriter(f)

		// 【生产模式】
		// 根据 gin 官方的说明：[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
		// 如果部署到生产环境，请使用以下模式：
		// 1.生产模式(release) 和开发模式的变化主要是禁用 gin 记录接口访问日志，
		// 2.go服务就必须使用nginx作为前置代理服务，这样也方便实现负载均衡
		// 3.如果程序发生 panic 等异常使用自定义的 panic 恢复中间件拦截、记录到日志

		// 生产模式，不使用 nginx
		lumberJackLogger := &lumberjack.Logger{ // 日志切割归档
			Filename:   variable.BasePath + variable.ConfigYml.GetString("Logs.GinLogName"), // 日志文件的位置
			MaxSize:    variable.ConfigYml.GetInt("Logs.MaxSize"),                           // 在进行切割之前，日志文件的最大大小（以MB为单位）
			MaxBackups: variable.ConfigYml.GetInt("Logs.MaxBackups"),                        // 保留旧文件的最大个数
			MaxAge:     variable.ConfigYml.GetInt("Logs.MaxAge"),                            // 保留旧文件的最大天数
			Compress:   variable.ConfigYml.GetBool("Logs.Compress"),                         // 是否压缩/归档旧文件
		}
		gin.DefaultWriter = io.MultiWriter(lumberJackLogger) // 设置gin 默认写入到lumberJackLogger
		router = gin_release.ReleaseRouter()
	} else {
		// 调试模式，开启 pprof 包，便于开发阶段分析程序性能
		router = gin.Default()
		pprof.Register(router)
	}
	// 设置可信任的代理服务器列表,gin (2021-11-24发布的v1.7.7版本之后出的新功能)
	if variable.ConfigYml.GetInt("HttpServer.TrustProxies.IsOpen") == 1 {
		if err := router.SetTrustedProxies(variable.ConfigYml.GetStringSlice("HttpServer.TrustProxies.ProxyServerList")); err != nil {
			variable.ZapLog.Error(consts.GinSetTrustProxyError, zap.Error(err))
		}
	} else {
		_ = router.SetTrustedProxies(nil)
	}

	// 根据配置进行设置跨域
	if variable.ConfigYml.GetBool("HttpServer.AllowCrossDomain") {
		router.Use(cors.Next())
	}

	router.GET("/", func(context *gin.Context) {
		context.String(http.StatusOK, "Api 模块接口 hello word！")
	})

	// 处理静态资源（不建议gin框架处理静态资源，参见 Public/readme.md 说明 ）
	router.Static("/public", "./public") //  定义静态资源路由与实际目录映射关系
	// router.StaticFile("/abcd", "./public/readme.md") // 可以根据文件名绑定需要返回的文件名

	// 总请求 API
	vApi := router.Group("/api/v1/")
	{
		// pod
		pods := vApi.Group("pods/")
		{
			// 第二个参数说明：
			// 1.它是一个表单参数验证器函数代码段，该函数从容器中解析，整个代码段略显复杂，但是对于使用者，您只需要了解用法即可，使用很简单，看下面 ↓↓↓
			// 2.编写该接口的验证器，位置：app/http/validator/api/pods/news.go
			// 3.将以上验证器注册在容器：app/http/validator/common/register_validator/api_register_validator.go  18 行为注册时的键（consts.ValidatorPrefix + "podList"）。那么获取的时候就用该键即可从容器获取
			pods.GET("list", validatorFactory.Create(consts.ValidatorPrefix+"PodList"))
			pods.GET("detail", validatorFactory.Create(consts.ValidatorPrefix+"Detail"))
			// pods.POST("", validatorFactory.Create(consts.ValidatorPrefix+"Delete")) //后端使用同一个接口，实现不同的请求
		}

		pod := vApi.Group("pod/")
		{
			pod.DELETE("del", validatorFactory.Create(consts.ValidatorPrefix+"Delete"))
			pod.PUT("update", validatorFactory.Create(consts.ValidatorPrefix+"Update"))
			pod.POST("create", validatorFactory.Create(consts.ValidatorPrefix+"Create"))

		}
	}
	return router
}
