package httpobs

import (
	"net/http"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

type SpanNameFunc func(r *http.Request) string

// TODO: option to collect request body and response, and skipper
func OpenTracing(fn SpanNameFunc) func(http.Handler) http.Handler {
	if fn == nil {
		fn = simpleSpanName
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Middleware Implementation for OpenTracing
			wireContext, _ := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))

			// Create the span referring to the RPC client if available.
			// If wireContext == nil, a root span will be created.
			serverSpan := opentracing.StartSpan(fn(r), ext.RPCServerOption(wireContext))
			defer serverSpan.Finish()

			ctx := opentracing.ContextWithSpan(r.Context(), serverSpan)
			r = r.WithContext(ctx)

			serverSpan.SetTag("request.method", r.Method)
			serverSpan.SetTag("request.url", r.URL.String())

			next.ServeHTTP(w, r)
		})
	}
}

func simpleSpanName(r *http.Request) string {
	return r.Method + " " + r.URL.Path
}
