package pool

import (
	"bytes"
	"sync"
)

type BufferPool struct {
	pool sync.Pool
}

func NewBufferPool(size int) *BufferPool {
	return &BufferPool{
		pool: sync.Pool{
			New: func() interface{} {
				return bytes.NewBuffer(make([]byte, 0, size))
			},
		},
	}
}

func (p *BufferPool) Get() *bytes.Buffer {
	return p.pool.Get().(*bytes.Buffer)
}

func (p *BufferPool) Put(b *bytes.Buffer) {
	b.Reset()
	p.pool.Put(b)
}
