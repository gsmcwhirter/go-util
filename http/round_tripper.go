package http

import (
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	"github.com/gsmcwhirter/go-util/v11/telemetry"
)

//counterfeiter:generate . RoundTripper
type RoundTripper = http.RoundTripper

type TelemeterRoundTripper struct {
	tel  *telemetry.Telemeter
	base http.RoundTripper

	spanOpts []telemetry.StartSpanOption
}

var _ http.RoundTripper = (*TelemeterRoundTripper)(nil)

func NewTelemeterRoundTripper(base http.RoundTripper, tel *telemetry.Telemeter, opts ...telemetry.StartSpanOption) *TelemeterRoundTripper {
	propagateRT := otelhttp.NewTransport(base)

	return &TelemeterRoundTripper{
		base:     propagateRT,
		tel:      tel,
		spanOpts: opts,
	}
}

func (rt *TelemeterRoundTripper) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	ctx, span := rt.tel.StartSpan(req.Context(), "http", "RoundTrip", rt.spanOpts...)
	defer span.End()

	defer func() {
		if resp != nil {
			span.SetAttributes(telemetry.HTTPAttributesFromHTTPStatusCode(resp.StatusCode)...)
		}

		if err != nil {
			span.SetStatus(telemetry.CodeError, err.Error())
		} else if resp != nil {
			code, reason := telemetry.SpanStatusFromHTTPStatusCode(resp.StatusCode)
			span.SetStatus(code, reason)
		}
	}()

	return rt.base.RoundTrip(req.WithContext(ctx))
}
