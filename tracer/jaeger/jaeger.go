package jaeger

import (
	"io"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/transport"
)

// NewTracer creates a jaeger tracer with constant sampler
// serviceURL should use this format => "http://jaeger:14268/api/traces"
//
func NewTracer(serviceName, serviceURL string) (opentracing.Tracer, io.Closer) {
	sampler := jaeger.NewConstSampler(true)
	reporter := jaeger.NewRemoteReporter(transport.NewHTTPTransport(serviceURL))
	return jaeger.NewTracer(serviceName, sampler, reporter)
}
