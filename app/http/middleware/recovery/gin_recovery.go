package recovery

import (
	"errors"
	"fmt"
	"time"

	"ginskeleton/app/global/consts"
	"ginskeleton/app/global/variable"
	"ginskeleton/app/utils/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// CustomRecovery 自定义错误(panic等)拦截中间件、对可能发生的错误进行拦截、统一记录
func CustomRecovery() gin.HandlerFunc {
	/*gin 还提供额外的Recovery中间件，CustomRecovery和RecoveryWithWriter
	Recovery()是另外俩个的封装，如果额外扩张，自定义ioWriter 或者handle可以用RecoveryWithWriter
	CustomRecovery 的DefaultErrorWriter 是将错误信息终端输出，如果要想写入自己的日志目录可以用RecoveryWithWriter，*/
	DefaultErrorWriter := &PanicExceptionRecord{}
	return gin.RecoveryWithWriter(DefaultErrorWriter, func(c *gin.Context, err interface{}) {
		// 这里针对发生的panic等异常进行统一响应即可
		// 这里的 err 数据类型为 ：runtime.boundsError  ，需要转为普通数据类型才可以输出
		response.ErrorSystem(c, "", fmt.Sprintf("%s", err))

		// 实现panic日志推送到邮件
		if variable.ConfigYml.GetInt("Email.IsToEmail") == 1 {
			if err := variable.EmailClient.SendMail(variable.ConfigYml.GetStringSlice("Email.To"),
				"GinWeb发生panic",
				fmt.Sprintf("%s----%s", time.Now().Format(variable.DateFormat), err)); err != nil {
				variable.ZapLog.Error("variable.EmailClient.SendMail error ", zap.Error(err))
			}
		}
	})
}

// PanicExceptionRecord  panic等异常记录
type PanicExceptionRecord struct{}

// 实现 write方法，RecoveryWithWriter需要
func (p *PanicExceptionRecord) Write(b []byte) (n int, err error) {
	errStr := string(b)
	err = errors.New(errStr)
	variable.ZapLog.Error(consts.ServerOccurredErrorMsg, zap.String("msg", errStr))
	return len(errStr), err
}
