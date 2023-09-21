package json

import (
	"io"

	sj "github.com/segmentio/encoding/json"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	"github.com/gsmcwhirter/go-util/v11/pool"
)

type Unmarshaler interface {
	UnmarshalJSON([]byte) error
}

func Unmarshal(b []byte, i interface{}) error {
	switch v := i.(type) {
	case Unmarshaler:
		return v.UnmarshalJSON(b)
	default:
		return sj.Unmarshal(b, i)
	}
}

var decPool = pool.NewBufferPool(4096)

func UnmarshalFromReader(r io.Reader, i interface{}) error {
	data := decPool.Get()
	defer decPool.Put(data)

	_, err := io.Copy(data, r)
	if err != nil {
		return err
	}

	return Unmarshal(data.Bytes(), i)
}

func ProtoUnmarshal(b []byte, m proto.Message) error {
	return ProtoUnmarshalOpts(b, m, protojson.UnmarshalOptions{})
}

func ProtoUnmarshalOpts(b []byte, m proto.Message, opts protojson.UnmarshalOptions) error {
	return opts.Unmarshal(b, m)
}

func ProtoUnmarshalFromReader(r io.Reader, m proto.Message) error {
	return ProtoUnmarshalFromReaderOpts(r, m, protojson.UnmarshalOptions{})
}

func ProtoUnmarshalFromReaderOpts(r io.Reader, m proto.Message, opts protojson.UnmarshalOptions) error {
	data := decPool.Get()
	defer decPool.Put(data)

	_, err := io.Copy(data, r)
	if err != nil {
		return err
	}

	return ProtoUnmarshalOpts(data.Bytes(), m, opts)
}
