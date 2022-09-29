package telemetry

import (
	"context"

	"go.opentelemetry.io/contrib/propagators/b3"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/metric/instrument/syncfloat64"
	"go.opentelemetry.io/otel/metric/instrument/syncint64"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdkTrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"go.opentelemetry.io/otel/trace"
)

type KeyValue = attribute.KeyValue

var (
	KVBool         = attribute.Bool
	KVBoolSlice    = attribute.BoolSlice
	KVInt          = attribute.Int
	KVIntSlice     = attribute.IntSlice
	KVInt64        = attribute.Int64
	KVInt64Slice   = attribute.Int64Slice
	KVFloat64      = attribute.Float64
	KVFloat64Slice = attribute.Float64Slice
	KVString       = attribute.String
	KVStringSlice  = attribute.StringSlice
	KVStringer     = attribute.Stringer
)

type Resource = resource.Resource

type (
	SpanExporter    = sdkTrace.SpanExporter
	Span            = trace.Span
	ReadOnlySpan    = trace.ReadOnlySpan
	StartSpanOption = trace.SpanStartOption
	TracerOption    = trace.TracerOption
	Tracer          = trace.Tracer
)

var WithAttributes = trace.WithAttributes

type (
	Meter         = metric.Meter
	MeterProvider = metric.MeterProvider
	MeterOption   = metric.MeterOption
)

type (
	Int64Counter       = syncint64.Counter
	Int64UpDownCounter = syncint64.UpDownCounter
	Int64Histogram     = syncint64.Histogram
)

type (
	Float64Counter       = syncfloat64.Counter
	Float64UpDownCounter = syncfloat64.UpDownCounter
	Float64Histogram     = syncfloat64.Histogram
)

type Telemeter struct { // trace.TracerProvider
	resource       *Resource
	tracerProvider *sdkTrace.TracerProvider
	traceExporter  SpanExporter
	meterProvider  MeterProvider
	propagator     propagation.TextMapPropagator
}

var (
	_ trace.TracerProvider = (*Telemeter)(nil)
	_ metric.MeterProvider = (*Telemeter)(nil)
)

func NewTelemeter(serviceName, version, instanceID string, spanExporter SpanExporter, meterProvider MeterProvider, spanSample float64) *Telemeter {
	res := newResource(serviceName, version, instanceID)
	return &Telemeter{
		resource:      res,
		traceExporter: spanExporter,
		meterProvider: meterProvider,
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

	global.SetMeterProvider(t)
}

func (t *Telemeter) Shutdown(ctx context.Context) error {
	return t.tracerProvider.Shutdown(ctx)
}

func (t *Telemeter) StartSpan(ctx context.Context, pkg, op string, opts ...StartSpanOption) (context.Context, Span) {
	return t.Tracer(pkg).Start(ctx, op, opts...)
}

func (t *Telemeter) Tracer(instrumentationName string, opts ...TracerOption) Tracer {
	return t.tracerProvider.Tracer(instrumentationName, opts...)
}

func (t *Telemeter) Meter(instrumentationName string, opts ...MeterOption) Meter {
	return t.meterProvider.Meter(instrumentationName, opts...)
}
