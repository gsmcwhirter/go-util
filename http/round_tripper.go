package http

import (
	"net/http"

	"github.com/gsmcwhirter/go-util/v10/telemetry"
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
	return &TelemeterRoundTripper{
		base:     base,
		tel:      tel,
		spanOpts: opts,
	}
}

func (rt *TelemeterRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	ctx, span := rt.tel.StartSpan(req.Context(), "http", "RoundTrip", rt.spanOpts...)
	defer span.End()

	return rt.base.RoundTrip(req.WithContext(ctx))
}
