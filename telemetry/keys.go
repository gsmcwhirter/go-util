package telemetry

import (
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
