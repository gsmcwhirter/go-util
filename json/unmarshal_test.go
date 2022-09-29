package json

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshal(t *testing.T) {
	t.Parallel()

	type testStruct struct {
		Foo string `json:"foo"`
		Bar int    `json:"baz"`
	}

	type args struct {
		b []byte
		i interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "1",
			args: args{
				b: []byte(`{"foo": "bar", "baz": 123}`),
				i: &testStruct{},
			},
			want: &testStruct{
				Foo: "bar",
				Bar: 123,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := Unmarshal(tt.args.b, tt.args.i)

			if (tt.wantErr && !assert.Error(t, err)) || (!tt.wantErr && !assert.NoError(t, err)) {
				return
			}
			if err != nil {
				return
			}

			assert.Equal(t, tt.want, tt.args.i)
		})
	}
}

func TestUnmarshalFromReader(t *testing.T) {
	t.Parallel()

	type testStruct struct {
		Foo string `json:"foo"`
		Bar int    `json:"baz"`
	}

	type args struct {
		r io.Reader
		i interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "1",
			args: args{
				r: bytes.NewBufferString(`{"foo": "bar", "baz": 123}`),
				i: &testStruct{},
			},
			want: &testStruct{
				Foo: "bar",
				Bar: 123,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := UnmarshalFromReader(tt.args.r, tt.args.i)

			if (tt.wantErr && !assert.Error(t, err)) || (!tt.wantErr && !assert.NoError(t, err)) {
				return
			}
			if err != nil {
				return
			}

			assert.Equal(t, tt.want, tt.args.i)
		})
	}
}
