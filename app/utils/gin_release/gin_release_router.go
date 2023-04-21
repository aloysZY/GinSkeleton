package gin_release

import (
	"ginskeleton/app/global/variable"
	"ginskeleton/app/http/middleware/access_log"
	"ginskeleton/app/http/middleware/context_timeout"
	"ginskeleton/app/http/middleware/recovery"
	"time"

	"github.com/gin-gonic/gin"
)

// ReleaseRouter 根据 gin 路由包官方的建议，gin 路由引擎如果在生产模式使用，官方建议设置为 release 模式
// 官方原版提示说明：[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
// 这里我们将按照官方指导进行生产模式精细化处理
func ReleaseRouter() *gin.Engine {
	// 切换到生产模式禁用 gin 输出接口访问日志，经过并发测试验证，可以提升5%的性能
	gin.SetMode(gin.ReleaseMode)
	// gin.DefaultWriter = ioutil.Discard   //生产模式，不使用 ngiinx 记录日志 注解掉这个行

	timeout := variable.ConfigYml.GetDuration("HttpServer.ContextTimeout")

	engine := gin.New()
	// 载入gin的中间件，关键是第二个中间件，我们对它进行了自定义重写，将可能的 panic 异常等，统一使用 zaplog 接管，保证全局日志打印统一
	engine.Use(access_log.GinAccessLogger(),
		recovery.CustomRecovery(),
		context_timeout.ContextTimeout(timeout*time.Second))
	return engine
}
