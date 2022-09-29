package json

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarshal(t *testing.T) {
	t.Parallel()

	type testStruct struct {
		Foo string `json:"foo"`
		Bar int    `json:"baz"`
	}

	type args struct {
		i interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "1",
			args: args{
				i: testStruct{
					Foo: "a",
					Bar: 123,
				},
			},
			want:    []byte(`{"foo":"a","baz":123}`),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := Marshal(tt.args.i)

			if (tt.wantErr && !assert.Error(t, err)) || (!tt.wantErr && !assert.NoError(t, err)) {
				return
			}
			if err != nil {
				return
			}

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestMarshalToBuffer(t *testing.T) {
	t.Parallel()

	type testStruct struct {
		Foo string `json:"foo"`
		Bar int    `json:"baz"`
	}

	type args struct {
		i interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "1",
			args: args{
				i: testStruct{
					Foo: "a",
					Bar: 123,
				},
			},
			want:    []byte(`{"foo":"a","baz":123}`),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := MarshalToBuffer(tt.args.i)

			if (tt.wantErr && !assert.Error(t, err)) || (!tt.wantErr && !assert.NoError(t, err)) {
				return
			}
			if err != nil {
				return
			}

			assert.Equal(t, tt.want, got)
		})
	}
}
