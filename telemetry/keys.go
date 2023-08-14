package telemetry

import (
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/semconv/v1.20.0/httpconv"
)

var (
	ClientRequest                       = httpconv.ClientRequest
	ClientRespone                       = httpconv.ClientResponse
	ServerRequest                       = httpconv.ServerRequest
	ClientStatus                        = httpconv.ClientStatus
	HTTPClientAttributesFromHTTPRequest = httpconv.ClientRequest
	HTTPServerAttributesFromHTTPRequest = httpconv.ServerRequest
	HTTPAttributesFromHTTPStatusCode    = httpconv.ClientResponse
	SpanStatusFromHTTPStatusCode        = httpconv.ClientStatus
)

func Handler(h http.Handler, name string, t *Telemeter) http.Handler {
	return otelhttp.NewHandler(h, name, otelhttp.WithTracerProvider(t), otelhttp.WithMeterProvider(t))
}
