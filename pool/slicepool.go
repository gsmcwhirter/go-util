package pool

import (
	"sync"
)

type SlicePool struct {
	pool sync.Pool
}

func NewSlicePool(size int) *SlicePool {
	return &SlicePool{
		pool: sync.Pool{
			New: func() interface{} {
				b := make([]byte, 0, size)
				return &b
			},
		},
	}
}

func (p *SlicePool) Get() []byte {
	b := p.pool.Get().(*[]byte)
	return *b
}

func (p *SlicePool) Put(b []byte) {
	b = b[:0]
	p.pool.Put(&b)
}
