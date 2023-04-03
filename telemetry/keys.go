package telemetry

import (
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
)

var (
	HTTPClientAttributesFromHTTPRequest = semconv.HTTPClientAttributesFromHTTPRequest
	HTTPServerAttributesFromHTTPRequest = semconv.HTTPServerAttributesFromHTTPRequest
	HTTPAttributesFromHTTPStatusCode    = semconv.HTTPAttributesFromHTTPStatusCode
	SpanStatusFromHTTPStatusCode        = semconv.SpanStatusFromHTTPStatusCode
)
