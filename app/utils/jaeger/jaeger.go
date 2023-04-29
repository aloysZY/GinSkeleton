package jaeger

import (
	"fmt"
	"io"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
)

func InitJaeger(serviceName, agentHost, agentPort string) (opentracing.Tracer, io.Closer, error) {
	// 根据配置初始化Tracer 返回Closer
	agentHostPort := fmt.Sprint(agentHost, ":", agentPort)
	tracer, closer, err := (&config.Configuration{
		ServiceName: serviceName,
		Sampler: &config.SamplerConfig{
			//Type: jaeger.SamplerTypeConst,
			Type: jaeger.SamplerTypeRemote,
			// param的值在0到1之间，设置为1则将所有的Operation输出到Reporter
			//Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans:            true,
			BufferFlushInterval: 1 * time.Second,
			LocalAgentHostPort:  agentHostPort,
		},
	}).NewTracer()
	if err != nil {
		return nil, nil, err
	}

	// 设置全局Tracer - 如果不设置将会导致上下文无法生成正确的Span
	opentracing.SetGlobalTracer(tracer)
	return tracer, closer, nil
}
