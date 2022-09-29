package json

import (
	"bytes"
	"testing"
)

var (
	Out []byte
	T   interface{}
)

func BenchmarkMarshal(b *testing.B) {
	type testStruct struct {
		Foo string `json:"foo"`
		Bar int    `json:"baz"`
	}

	t := testStruct{
		Foo: "abc",
		Bar: 123,
	}

	for i := 0; i < b.N; i++ {
		out, err := Marshal(t)
		Out = out

		if err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkMarshalToBuffer(b *testing.B) {
	type testStruct struct {
		Foo string `json:"foo"`
		Bar int    `json:"baz"`
	}

	t := testStruct{
		Foo: "abc",
		Bar: 123,
	}

	for i := 0; i < b.N; i++ {
		out, err := MarshalToBuffer(t)
		Out = out

		if err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkUnmarshal(b *testing.B) {
	type testStruct struct {
		Foo string `json:"foo"`
		Bar int    `json:"baz"`
	}

	byt := []byte(`{"foo": "bar", "baz": 123}`)
	t := testStruct{}

	for i := 0; i < b.N; i++ {
		err := Unmarshal(byt, &t)
		if err != nil {
			b.Error(err)
		}

		T = t
	}
}

func BenchmarkUnmarshalFromReader(b *testing.B) {
	type testStruct struct {
		Foo string `json:"foo"`
		Bar int    `json:"baz"`
	}

	byt := `{"foo": "bar", "baz": 123}`
	t := testStruct{}

	for i := 0; i < b.N; i++ {
		buf := bytes.NewBufferString(byt)

		err := UnmarshalFromReader(buf, &t)
		if err != nil {
			b.Error(err)
		}

		T = t
	}
}
