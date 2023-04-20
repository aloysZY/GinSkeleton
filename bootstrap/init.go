package bootstrap

import (
	"fmt"
	"ginskeleton/app/utils/kube_client"
	"log"
	"os"
	"path/filepath"
	"time"

	_ "ginskeleton/app/core/destroy" // 监听程序退出信号，用于资源的释放
	"ginskeleton/app/global/my_errors"
	"ginskeleton/app/global/variable"
	"ginskeleton/app/http/validator/common/register_validator"
	"ginskeleton/app/service/sys_log_hook"
	"ginskeleton/app/utils/casbin_v2"
	"ginskeleton/app/utils/email"
	"ginskeleton/app/utils/gorm_v2"
	"ginskeleton/app/utils/snow_flake"
	"ginskeleton/app/utils/validator_translation"
	"ginskeleton/app/utils/websocket/core"
	"ginskeleton/app/utils/yml_config"
	"ginskeleton/app/utils/zap_factory"
)

// 检查项目必须的非编译目录是否存在，避免编译后调用的时候缺失相关目录
func checkRequiredFolders() {
	// 1.设置默认配置文件路径
	if _, err := os.Stat(variable.BasePath + "/config/config.yml"); err != nil {
		log.Fatal(my_errors.ErrorsConfigYamlNotExists + err.Error())
	}
	if _, err := os.Stat(variable.BasePath + "/config/gorm_v2.yml"); err != nil {
		log.Fatal(my_errors.ErrorsConfigGormNotExists + err.Error())
	}

	// 2.检查public目录是否存在
	if _, err := os.Stat(variable.BasePath + "/public"); err != nil {
		/*log.Fatal(my_errors.ErrorsPublicNotExists + err.Error())*/
		log.Println(my_errors.ErrorsPublicNotExists + "：beginning Create directory")
		// 创建目录的时候，最后要有"/"
		err := os.MkdirAll(filepath.Dir(variable.BasePath+"/public/"), 0755)
		if err != nil && !os.IsExist(err) {
			log.Fatal(my_errors.ErrorCreateFailed + err.Error())
		}
	}

	// 3.检查storage/logs 目录是否存在
	if _, err := os.Stat(variable.BasePath + "/storage/logs"); err != nil {
		log.Println(my_errors.ErrorsStorageLogsNotExists + "：beginning Create directory")
		err := os.MkdirAll(filepath.Dir(variable.BasePath+"/storage/logs/"), 0755)
		if err != nil && !os.IsExist(err) {
			log.Fatal(my_errors.ErrorCreateFailed + err.Error())
		}
	}

	// 4.自动创建软连接、更好的管理静态资源，上传的文件链接到public目录，nginx配置也是这个，
	if _, err := os.Stat(variable.BasePath + "/public"); err == nil {
		if err = os.RemoveAll(variable.BasePath + "/public"); err != nil {
			log.Fatal(my_errors.ErrorsSoftLinkDeleteFail + err.Error())
		}
	}

	// 5.创建/storage/app 这是文件上传目录
	if _, err := os.Stat(variable.BasePath + "/storage/app"); err != nil {
		log.Println(my_errors.ErrorsStorageLogsNotExists + "：beginning Create directory")
		if err := os.MkdirAll(filepath.Dir(variable.BasePath+"/storage/app/"), 0755); err != nil && !os.IsExist(err) {
			log.Fatal(my_errors.ErrorCreateFailed + err.Error())
		}
	}

	// 6./storage/app 这是文件上传目录创建软连接，public/storage
	if err := os.Symlink(variable.BasePath+"/storage/app", variable.BasePath+"/public"); err != nil {
		log.Fatal(my_errors.ErrorsSoftLinkCreateFail + err.Error())
	}

}

func init() {
	// 1. 初始化 项目根路径，参见 variable 常量包，相关路径：app\global\variable\variable.go

	// 2.检查配置文件以及日志目录等非编译性的必要条件
	checkRequiredFolders()
	log.Println("Check required succeeded!")

	// 3.初始化表单参数验证器，注册在容器（Web、Api共用容器）
	// register_validator.WebRegisterValidator()
	register_validator.ApiRegisterValidator()
	log.Println("表单参数验证器 succeeded!")

	// 4.启动针对配置文件(confgi.yml、gorm_v2.yml)变化的监听， 配置文件操作指针，初始化为全局变量
	variable.ConfigYml = yml_config.CreateYamlFactory()
	variable.ConfigYml.ConfigFileChangeListen()
	log.Println("应用配置文件加载成功!")
	// config>gorm_v2.yml 启动文件变化监听事件
	variable.ConfigGormv2Yml = variable.ConfigYml.Clone("gorm_v2")
	variable.ConfigGormv2Yml.ConfigFileChangeListen()
	log.Println("数据库配置文件加载成功!")

	// 5.初始化全局日志句柄，并载入日志钩子处理函数
	variable.ZapLog = zap_factory.CreateZapFactory(sys_log_hook.ZapLogHandler)
	variable.ZapLog.Info("日志初始化完成！")

	// 6.根据配置初始化 gorm mysql 全局 *gorm.Db
	if variable.ConfigGormv2Yml.GetInt("Gormv2.Mysql.IsInitGlobalGormMysql") == 1 {
		if dbMysql, err := gorm_v2.GetOneMysqlClient(); err != nil {
			log.Fatal(my_errors.ErrorsGormInitFail + err.Error())
		} else {
			variable.ZapLog.Info("mysql数据库初始化完成！")
			variable.GormDbMysql = dbMysql
		}
	}
	// 根据配置初始化 gorm sqlserver 全局 *gorm.Db
	if variable.ConfigGormv2Yml.GetInt("Gormv2.Sqlserver.IsInitGlobalGormSqlserver") == 1 {
		if dbSqlserver, err := gorm_v2.GetOneSqlserverClient(); err != nil {
			log.Fatal(my_errors.ErrorsGormInitFail + err.Error())
		} else {
			variable.ZapLog.Info("sqlserver数据库初始化完成！")
			variable.GormDbSqlserver = dbSqlserver
		}
	}
	// 根据配置初始化 gorm postgresql 全局 *gorm.Db
	if variable.ConfigGormv2Yml.GetInt("Gormv2.PostgreSql.IsInitGlobalGormPostgreSql") == 1 {
		if dbPostgre, err := gorm_v2.GetOnePostgreSqlClient(); err != nil {
			log.Fatal(my_errors.ErrorsGormInitFail + err.Error())
		} else {
			variable.ZapLog.Info("postgresql数据库初始化完成！")
			variable.GormDbPostgreSql = dbPostgre
		}
	}

	// 7.雪花算法全局变量
	variable.SnowFlake = snow_flake.CreateSnowflakeFactory()
	variable.ZapLog.Info("雪花算法初始化完成！")

	// 8.websocket Hub中心启动
	if variable.ConfigYml.GetInt("Websocket.Start") == 1 {
		// websocket 管理中心hub全局初始化一份
		variable.WebsocketHub = core.CreateHubFactory()
		if Wh, ok := variable.WebsocketHub.(*core.Hub); ok {
			go Wh.Run()
		}
	}

	// 9.casbin 依据配置文件设置参数(IsInit=1)初始化
	if variable.ConfigYml.GetInt("Casbin.IsInit") == 1 {
		var err error
		if variable.Enforcer, err = casbin_v2.InitCasbinEnforcer(); err != nil {
			log.Fatal(err.Error())
		}
	}

	// 10.全局注册 validator 错误翻译器,zh 代表中文，en 代表英语
	if err := validator_translation.InitTrans("zh"); err != nil {
		log.Fatal(my_errors.ErrorsValidatorTransInitFail + err.Error())
	}

	// 11.初始化clientset
	/*	if variable.ConfigYml.GetInt("Kubernetes.Clientset.IsInitGlobalClientset") == 1 {
		if client, err := kube_client.GetClientset(variable.ConfigYml.GetString("Kubernetes.Clientset.ConfigPath")); err != nil {
			variable.ZapLog.Fatal("Kubernetes.Clientset error: ", zap.Error(err))
		} else {
			variable.ZapLog.Info("kubernetes client success!")
			variable.Clientset = client.Clientset
		}
	}*/
	//if variable.ConfigYml.GetInt("Kubernetes.Clientset.IsInitGlobalClientset") == 1 {
	//	if client, err := kube_client.GetOneClientset(variable.ConfigYml.GetString("Kubernetes.Clientset.ConfigPath")); err != nil {
	//		variable.ZapLog.Fatal("Kubernetes.Clientset error: ", zap.Error(err))
	//	} else {
	//		variable.ZapLog.Info("kubernetes client success!")
	//		variable.Clientset = client
	//	}
	//}
	if variable.ConfigYml.GetInt("Kubernetes.Clientset.IsInitGlobalClientset") == 1 {
		variable.KubeControllerClientset = kube_client.NewKubeControllerclientset(variable.ConfigYml.GetString("Kubernetes.Clientset.ConfigPath"), 0)
		stopPodch := make(chan struct{})
		go func() {
			variable.KubeControllerClientset.Run(stopPodch)
			<-stopPodch
		}()
		for {
			if variable.KubeControllerClientset.Status == 1 {
				break
			}
			time.Sleep(time.Second * 1)
			fmt.Println("sleep 1S")
		}
	}

	// 12.如果设置发送邮件，就初始化邮件
	if variable.ConfigYml.GetInt("Email.IsToEmail") == 1 {
		variable.EmailClient = email.NewEmail()
	}
}
