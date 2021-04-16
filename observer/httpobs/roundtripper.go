package httpobs

import (
	"net/http"

	"github.com/opentracing/opentracing-go"
)

// RoundTripper create http.RoundTripper that will propagate headers
func RoundTripper(original http.RoundTripper) http.RoundTripper {
	return &transport{original: original}
}

type transport struct {
	original http.RoundTripper
}

func (rt *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	if rt.original == nil {
		rt.original = http.DefaultTransport
	}

	span, _ := opentracing.StartSpanFromContext(req.Context(), req.Method+" "+req.URL.Path)
	defer span.Finish()
	span.SetTag("http.method", req.Method)
	span.SetTag("http.url", req.URL.String())
	err := opentracing.GlobalTracer().Inject(span.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(req.Header))
	if err != nil {
		return rt.original.RoundTrip(req)
	}

	res, err := rt.original.RoundTrip(req)
	if res != nil {
		span.SetTag("http.status_code", res.StatusCode)
	}
	return res, err
}
