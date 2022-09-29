package request

import (
	"context"
	"testing"
)

func TestNewRequestContext(t *testing.T) {
	t.Parallel()

	t.Run("check request id inserted", func(t *testing.T) {
		t.Parallel()
		got := NewRequestContext()
		if rid, ok := got.Value(ContextKey("request_id")).(string); !ok {
			t.Error("no request_id was present")
		} else if rid == "" {
			t.Error("request_id was present but empty")
		}
	})
}

func TestNewRequestContextFrom(t *testing.T) {
	t.Parallel()

	type args struct {
		origID string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			"no origID",
			args{},
		},
		{
			"with origID",
			args{"test"},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			if tt.args.origID != "" {
				ctx = NewRequestContextWithRequestID(ctx, tt.args.origID)
				if rid, ok := ctx.Value(ContextKey("request_id")).(string); !ok {
					t.Error("no origID request_id was present")
				} else if rid == "" {
					t.Error("origID request_id was present but empty")
				}
			}

			got := NewRequestContextFrom(ctx)
			if rid, ok := got.Value(ContextKey("request_id")).(string); !ok {
				t.Error("no request_id was present")
			} else if rid == "" {
				t.Error("request_id was present but empty")
			} else if tt.args.origID != "" && rid == tt.args.origID {
				t.Error("request_id was the same as the origID")
			}
		})
	}
}

func TestGetRequestID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		origID string
		want   string
		wantOk bool
	}{
		{
			"none",
			"",
			"",
			false,
		},
		{
			"some",
			"test",
			"test",
			true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			if tt.origID != "" {
				ctx = NewRequestContextWithRequestID(ctx, tt.origID)
				if rid, ok := ctx.Value(ContextKey("request_id")).(string); !ok {
					t.Error("no origID request_id was present")
				} else if rid == "" {
					t.Error("origID request_id was present but empty")
				}
			}

			got, ok := GetRequestID(ctx)
			if ok != tt.wantOk {
				t.Errorf("request_id ok got %v want %v", ok, tt.wantOk)
			}

			if got != tt.want {
				t.Errorf("request_id got %v want %v", got, tt.want)
			}
		})
	}
}
