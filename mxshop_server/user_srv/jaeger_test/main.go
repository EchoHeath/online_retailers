package main

import (
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"time"
)

func main() {
	cfg := &config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans: true,
			LocalAgentHostPort:"127.0.0.1:6831",
		},
		ServiceName: "mxshop",
	}

	tracer, closer, err := cfg.NewTracer(config.Logger(jaeger.StdLogger))
	if err != nil {
		panic(err)
	}

	defer closer.Close()

	parentSpan := tracer.StartSpan("main")
	span := tracer.StartSpan("funcA", opentracing.ChildOf(parentSpan.Context()))
	time.Sleep(time.Millisecond * 500)
	span.Finish()

	span2 := tracer.StartSpan("funcA", opentracing.ChildOf(span.Context()))
	time.Sleep(time.Millisecond * 1000)
	span2.Finish()

	parentSpan.Finish()
	//tracer, closer, err := cfg.New(service, config.Logger(jaeger.StdLogger))
	//if err != nil {
	//	panic(fmt.Sprintf("ERROR: cannot init Jaeger: %v\n", err))
	//}
}

