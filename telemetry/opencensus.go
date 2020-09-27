package telemetry

import (
	"context"

	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
	"go.opencensus.io/trace"

	"github.com/gsmcwhirter/go-util/v7/errors"
	"github.com/gsmcwhirter/go-util/v7/request"
)

type View = view.View
type ViewData = view.Data

var CountView = view.Count
var RegisterView = view.Register

type Measurement = stats.Measurement

var Int64 = stats.Int64

type TagKey = tag.Key

var NewTagKey = tag.NewKey
var MustNewTagKey = tag.MustNewKey

type SpanData = trace.SpanData

type TraceExporter = trace.Exporter
type ViewExporter = view.Exporter

var StringAttribute = trace.StringAttribute

// ErrNoMeasurements is the error from Census.Record when no measurements were provided
var ErrNoMeasurements = errors.New("no measurements provided")

type Tag struct {
	Key TagKey
	Val string
}

type Census struct {
	statsExporter view.Exporter
	statsFlush    func(view.Exporter)
	traceExporter trace.Exporter
	traceFlush    func(trace.Exporter)
}

type Options struct {
	StatsExporter      view.Exporter
	StatsExporterFlush func(view.Exporter)
	TraceExporter      trace.Exporter
	TraceExporterFlush func(trace.Exporter)
	TraceProbability   float64
}

func NewCensus(opts Options) *Census {
	c := &Census{
		statsExporter: opts.StatsExporter,
		statsFlush:    opts.StatsExporterFlush,
		traceExporter: opts.TraceExporter,
		traceFlush:    opts.TraceExporterFlush,
	}

	view.RegisterExporter(c.statsExporter)
	trace.RegisterExporter(c.traceExporter)

	trace.ApplyConfig(trace.Config{DefaultSampler: trace.ProbabilitySampler(opts.TraceProbability)})

	return c
}

func (c *Census) Flush() {
	if c.statsFlush != nil {
		c.statsFlush(c.statsExporter)
	}

	if c.traceFlush != nil {
		c.traceFlush(c.traceExporter)
	}
}

func (c *Census) startSpan(ctx context.Context, name string, kind trace.StartOption, keyvals ...string) (context.Context, *trace.Span) {
	ctx, span := trace.StartSpan(ctx, name, kind)

	attributes := make([]trace.Attribute, 0, len(keyvals)/2+1)

	if rid, ok := request.GetRequestID(ctx); ok {
		attributes = append(attributes, trace.StringAttribute("request_id", rid))
	}

	for i := 0; i < len(keyvals)-1; i += 2 {
		attributes = append(attributes, trace.StringAttribute(keyvals[i], keyvals[i+1]))
	}

	if len(attributes) > 0 {
		span.AddAttributes(attributes...)
	}

	return ctx, span
}

func (c *Census) StartSpan(ctx context.Context, name string, keyvals ...string) (context.Context, *trace.Span) {
	return c.startSpan(ctx, name, trace.WithSpanKind(trace.SpanKindServer), keyvals...)
}

func (c *Census) StartClientSpan(ctx context.Context, name string, keyvals ...string) (context.Context, *trace.Span) {
	return c.startSpan(ctx, name, trace.WithSpanKind(trace.SpanKindServer), keyvals...)
}

func (c *Census) AddSpanAttributes(span *trace.Span, keyvals ...string) {
	attributes := make([]trace.Attribute, 0, len(keyvals)/2)

	for i := 0; i < len(keyvals)-1; i += 2 {
		attributes = append(attributes, trace.StringAttribute(keyvals[i], keyvals[i+1]))
	}

	if len(attributes) > 0 {
		span.AddAttributes(attributes...)
	}
}

func (c *Census) Record(ctx context.Context, ms []Measurement, tags ...Tag) error {
	if len(ms) == 0 {
		return ErrNoMeasurements
	}

	muts := make([]tag.Mutator, 0, len(tags))
	for _, t := range tags {
		muts = append(muts, tag.Upsert(t.Key, t.Val))
	}

	opts := []stats.Options{
		stats.WithMeasurements(ms...),
		stats.WithTags(muts...),
	}

	return stats.RecordWithOptions(ctx, opts...)
}
