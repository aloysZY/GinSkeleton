package tracering

import (
	"context"

	"ginskeleton/app/global/variable"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/uber/jaeger-client-go"
)

func Tracering() func(c *gin.Context) {
	return func(c *gin.Context) {
		var newCtx context.Context
		var span opentracing.Span

		spanCtx, err := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(c.Request.Header))
		if err != nil {
			span, newCtx = opentracing.StartSpanFromContextWithTracer(c.Request.Context(), variable.Tracer, c.Request.URL.Path)
		} else {
			span, newCtx = opentracing.StartSpanFromContextWithTracer(
				c.Request.Context(),
				variable.Tracer,
				c.Request.URL.Path,
				opentracing.ChildOf(spanCtx),
				opentracing.Tag{Key: string(ext.Component), Value: "HTTP"},
			)
		}
		defer span.Finish()

		spacnDb, _ := opentracing.StartSpanFromContextWithTracer(
			c.Request.Context(),
			variable.Tracer,
			c.Request.URL.Path,
		)
		ctx := opentracing.ContextWithSpan(newCtx, spacnDb)

		// 6. 将上下文传入DB实例，生成Session会话
		// 这样子就能把这个会话的全部信息反馈给Jaeger
		variable.GormDbMysql = variable.GormDbMysql.WithContext(ctx)

		var traceID string
		var spanID string
		var spanContextID = span.Context()
		switch spanContextID.(type) {
		case jaeger.SpanContext:
			jaegerContextID := spanContextID.(jaeger.SpanContext)
			traceID = jaegerContextID.TraceID().String()
			spanID = jaegerContextID.SpanID().String()
		}
		c.Set("X-Trace-ID", traceID)
		c.Set("X-Span-ID", spanID)

		c.Request = c.Request.WithContext(newCtx)
		c.Next()
	}
}