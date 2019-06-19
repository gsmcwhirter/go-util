package census

import (
	"context"

	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
	"go.opencensus.io/trace"

	"github.com/gsmcwhirter/go-util/v4/errors"
	"github.com/gsmcwhirter/go-util/v4/logging"
)

type View = view.View

var CountView = view.Count
var RegisterView = view.Register

var Int64 = stats.Int64

type TagKey = tag.Key

var NewTagKey = tag.NewKey
var NewTag = tag.New
var InsertTag = tag.Insert

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

func (c *OpenCensus) StartSpan(ctx context.Context, name string) (context.Context, *trace.Span) {
	return trace.StartSpan(ctx, name)
}

func (c *OpenCensus) Record(ctx context.Context, ms ...stats.Measurement) {
	stats.Record(ctx, ms...)
}
