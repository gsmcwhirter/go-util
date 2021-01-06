package json

import (
	sj "github.com/segmentio/encoding/json"

	"github.com/gsmcwhirter/go-util/v8/pool"
)

type Marshaler interface {
	MarshalJSON() ([]byte, error)
}

var encPool = pool.NewSlicePool(4096)

func Marshal(i interface{}) ([]byte, error) {
	switch v := i.(type) {
	case Marshaler:
		return v.MarshalJSON()
	default:
		return MarshalToBuffer(i)
	}
}

func MarshalToBuffer(i interface{}) ([]byte, error) {
	var err error

	b := encPool.Get()
	defer encPool.Put(b)

	if b.Data, err = sj.Append(b.Data[:0], i, 0); err != nil {
		return nil, err
	}

	d := make([]byte, len(b.Data))
	copy(d, b.Data)
	return d, nil
}
