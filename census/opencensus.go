package census

import (
	"context"

	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
	"go.opencensus.io/trace"

	"github.com/gsmcwhirter/go-util/v4/errors"
	"github.com/gsmcwhirter/go-util/v4/logging"
	"github.com/gsmcwhirter/go-util/v4/request"
)

type View = view.View
type ViewData = view.Data

var CountView = view.Count
var RegisterView = view.Register

var Int64 = stats.Int64

type TagKey = tag.Key

var NewTagKey = tag.NewKey
var NewTag = tag.New
var InsertTag = tag.Insert

type SpanData = trace.SpanData

var ErrBadExporter = errors.New("unsupported exporter")

type dependencies interface {
	Logger() logging.Logger
}

type OpenCensus struct {
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

func NewCensus(deps dependencies, opts Options) *OpenCensus {
	c := &OpenCensus{
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

func (c *OpenCensus) Flush() {
	if c.statsFlush != nil {
		c.statsFlush(c.statsExporter)
	}

	if c.traceFlush != nil {
		c.traceFlush(c.traceExporter)
	}
}

func (c *OpenCensus) StartSpan(ctx context.Context, name string, keyvals ...string) (context.Context, *trace.Span) {
	ctx, span := trace.StartSpan(ctx, name)

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

func (c *OpenCensus) Record(ctx context.Context, ms ...stats.Measurement) {
	stats.Record(ctx, ms...)
}
