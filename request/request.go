package request

import (
	"context"

	"github.com/rs/xid"
)

// ContextKey is a wrapper type for our keys attached to a context
type ContextKey string

// NewRequestContext creates a new request context from context.Background()
func NewRequestContext() context.Context {
	return NewRequestContextFrom(context.Background())
}

// NewRequestContextFrom creates a new request context from an existing context
// regenerating a request_id value
func NewRequestContextFrom(ctx context.Context) context.Context {
	return NewRequestContextWithRequestID(ctx, GenerateRequestID())
}

// NewRequestContextWithRequestID creates a new request context from an existing context
// with the provided request id rid
func NewRequestContextWithRequestID(ctx context.Context, rid string) context.Context {
	return context.WithValue(ctx, ContextKey("request_id"), rid)
}

// GenerateRequestID generates a new random-enough request id for a request context
func GenerateRequestID() string {
	return xid.New().String()
}

func GetRequestID(ctx context.Context) (string, bool) {
	rid, ok := ctx.Value(ContextKey("request_id")).(string)
	return rid, ok
}
