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
		var ctx context.Context
		var span opentracing.Span

		// 设置中间件
		spanCtx, err := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(c.Request.Header))
		if err != nil {
			span, ctx = opentracing.StartSpanFromContextWithTracer(c.Request.Context(), variable.Tracer, c.Request.URL.Path)
		} else {
			span, ctx = opentracing.StartSpanFromContextWithTracer(
				c.Request.Context(),
				variable.Tracer,
				c.Request.URL.Path,
				opentracing.ChildOf(spanCtx),
				opentracing.Tag{Key: string(ext.Component), Value: "HTTP"},
			)
		}

		defer span.Finish()

		// 记录日志使用用的 ID 信息
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

		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
