package pool

import (
	"sync"
)

type ByteSlice struct {
	Data []byte
}

type SlicePool struct {
	pool sync.Pool
}

func NewSlicePool(size int) *SlicePool {
	return &SlicePool{
		pool: sync.Pool{
			New: func() interface{} {
				return &ByteSlice{Data: make([]byte, 0, size)}
			},
		},
	}
}

func (p *SlicePool) Get() *ByteSlice {
	return p.pool.Get().(*ByteSlice)
}

func (p *SlicePool) Put(b *ByteSlice) {
	b.Data = b.Data[:0]
	p.pool.Put(b)
}
