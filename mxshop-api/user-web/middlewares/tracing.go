package middlewares

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"mxshop-api/user-web/global"
)

func Trace() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		cfg := &config.Configuration{
			Sampler: &config.SamplerConfig{
				Type:  jaeger.SamplerTypeConst,
				Param: 1,
			},
			Reporter: &config.ReporterConfig{
				LogSpans: true,
				LocalAgentHostPort: fmt.Sprintf("%s:%d", global.ServerConfig.JaegerInfo.Host, global.ServerConfig.JaegerInfo.Port),
			},
			ServiceName: global.ServerConfig.JaegerInfo.Name,
		}

		tracer, closer, err := cfg.NewTracer(config.Logger(jaeger.StdLogger))
		if err != nil {
			panic(err)
		}

		defer closer.Close()

		startSpan := tracer.StartSpan(ctx.Request.URL.Path)
		defer startSpan.Finish()

		ctx.Set("trace", tracer)
		ctx.Set("parentSpan", startSpan)
		ctx.Next()
	}
}
