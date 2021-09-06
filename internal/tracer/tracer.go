package tracer

import (
	"github.com/opentracing/opentracing-go"
	"github.com/ozonva/ova-journey-api/internal/config"
	"github.com/rs/zerolog/log"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	"github.com/uber/jaeger-lib/metrics"
	"io"
)

// InitTracer - initialize Jaeger tracer
func InitTracer(serviceName string, endpoint *config.EndpointConfiguration) io.Closer {
	// Sample configuration for testing. Use constant sampling to sample every trace
	// and enable LogSpan to log every span via configured Logger.
	cfg := jaegercfg.Configuration{
		ServiceName: serviceName,
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: endpoint.GetEndpointAddress(),
		},
	}

	jLogger := jaegerlog.StdLogger
	jMetricsFactory := metrics.NullFactory

	// Initialize tracer with a logger and a metrics factory
	tracer, closer, err := cfg.NewTracer(jaegercfg.Logger(jLogger), jaegercfg.Metrics(jMetricsFactory))
	if err != nil {
		log.Error().Err(err).Msg("Could not initialize jaeger tracer")
		return nil
	}

	opentracing.SetGlobalTracer(tracer)

	return closer
}
