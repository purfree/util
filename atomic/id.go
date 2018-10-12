package atomic

import "sync/atomic"

type ID struct {
	id uint64
}

func (p *ID) Add() uint64 {
	id := atomic.AddUint64(&p.id, 1)
	return id
}
