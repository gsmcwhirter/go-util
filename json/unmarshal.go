package json

import (
	"io"

	sj "github.com/segmentio/encoding/json"

	"github.com/gsmcwhirter/go-util/v8/pool"
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
