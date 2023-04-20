package access_log

import (
	"fmt"
	"time"

	"ginskeleton/app/global/variable"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

// GinAccessLogger 接收gin框架默认的日志，添加访问日志记录
func GinAccessLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		// query := c.Request.URL.RawQuery
		/*		// 读取 body 内容,数据直接读取到缓存，数据太大会内存消耗严重,上传文件的时候内内占用大,等待时间长
				bodyByte, _ := ioutil.ReadAll(c.Request.Body)
				// 将读取的内容重新赋值，不然上面读取后之后的路由不能读取
				c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyByte))*/
		c.Next()
		cost := time.Since(start)
		/*		// 读取 body 内容,数据直接读取到缓存，数据太大会内存消耗严重,上传文件的时候内内占用大,等待时间长
				bodyByte, _ := ioutil.ReadAll(c.Request.Body)
				// 将读取的内容重新赋值，不然上面读取后之后的路由不能读取
				c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyByte))*/
		// 如果太大了就不写入日志,后面分别判断 body 大小，然后如果小于一定值记录，不然记录空？
		// if c.Request.Body.Read()
		accessStr := fmt.Sprintf(`{"time":"%v","status":"%d","method":"%s","path":"%s",query:"%s","user-agent":"%s","cost":"%v"}`, start.Format(variable.DateFormat), c.Writer.Status(), c.Request.Method, c.Request.URL.Path, c.Request.URL.RawQuery, c.Request.UserAgent(), cost)
		// 输出到gin.DefaultWriter
		if _, err := fmt.Fprintln(gin.DefaultWriter, accessStr); err != nil {
			variable.ZapLog.Error("fmt.Fprintln(gin.DefaultWriter, accessStr) error ", zap.Error(err))
		}
	}
}
