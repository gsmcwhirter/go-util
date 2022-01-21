package telemetry

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"

	"github.com/gsmcwhirter/go-util/v9/json"
)

func TelemTest1(ctx context.Context, t *Telemeter) {
	ctx, span := t.StartSpan(ctx, "test", "TelemTest1")
	defer span.End()

	time.Sleep(5 * time.Millisecond)
	TelemTest2(ctx, t)
	time.Sleep(5 * time.Millisecond)
}

func TelemTest2(ctx context.Context, t *Telemeter) {
	_, span := t.StartSpan(ctx, "test", "TelemTest2")
	defer span.End()

	time.Sleep(10 * time.Millisecond)
}

type span struct {
	Name string `json:"Name"`
}

func TestTelemetry(t *testing.T) {
	w := &bytes.Buffer{}
	w.WriteString("[")
	exp, err := stdouttrace.New(
		stdouttrace.WithWriter(w),
		// stdouttrace.WithPrettyPrint(),
		stdouttrace.WithoutTimestamps(),
	)
	assert.NoError(t, err, "failed to construct stdout exporter")

	tm := NewTelemeter("test", "v0", "test_instance", exp, 1.0)
	ctx := context.Background()

	TelemTest1(ctx, tm)

	err = tm.Shutdown(ctx)
	assert.NoError(t, err, "telemeter shutdown error")
	w.WriteString("]")

	out := strings.ReplaceAll(w.String(), "}\n{", "},\n{")

	fmt.Println(out)

	var spans []span
	err = json.Unmarshal([]byte(out), &spans)
	assert.NoError(t, err, "unmarshal failed")

	assert.Equal(t, []span{
		{Name: "TelemTest2"},
		{Name: "TelemTest1"},
	}, spans)
}
