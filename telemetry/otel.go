package telemetry

import (
	"context"

	"go.opentelemetry.io/contrib/propagators/b3"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdkTrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"go.opentelemetry.io/otel/trace"
)

type KeyValue = attribute.KeyValue
type SpanExporter = sdkTrace.SpanExporter
type Resource = resource.Resource
type Span = trace.Span
type StartSpanOption = trace.SpanStartOption
type TracerOption = trace.TracerOption
type Tracer = trace.Tracer

type Telemeter struct { // trace.TracerProvider
	resource       *Resource
	tracerProvider *sdkTrace.TracerProvider
	traceExporter  SpanExporter
	propagator     propagation.TextMapPropagator
}

var _ trace.TracerProvider = (*Telemeter)(nil)

func NewTelemeter(serviceName, version, instanceID string, spanExporter SpanExporter, spanSample float64) *Telemeter {
	res := newResource(serviceName, version, instanceID)
	return &Telemeter{
		resource:      res,
		traceExporter: spanExporter,
		tracerProvider: sdkTrace.NewTracerProvider(
			sdkTrace.WithBatcher(spanExporter),
			sdkTrace.WithResource(res),
			sdkTrace.WithSampler(
				sdkTrace.ParentBased(
					sdkTrace.TraceIDRatioBased(spanSample),
				),
			),
		),
		propagator: propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			b3.New(),
		),
	}
}

func newResource(serviceName, version, instanceID string) *Resource {
	return resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(serviceName),
		semconv.ServiceVersionKey.String(version),
		semconv.ServiceInstanceIDKey.String(instanceID),
	)
}

func (t *Telemeter) MakeDefault() {
	otel.SetTracerProvider(t)
	otel.SetTextMapPropagator(t.propagator)
}

func (t *Telemeter) Shutdown(ctx context.Context) error {
	return t.tracerProvider.Shutdown(ctx)
}

func (t *Telemeter) StartSpan(ctx context.Context, pkg, op string, opts ...StartSpanOption) (context.Context, Span) {
	return t.Tracer(pkg).Start(ctx, op, opts...)
}

func (t *Telemeter) Tracer(instrumentationName string, opts ...TracerOption) Tracer {
	return t.tracerProvider.Tracer(instrumentationName)
}
