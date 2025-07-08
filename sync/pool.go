package sync

import "sync"

type Pool[V any] struct {
	pool sync.Pool
	New  func() V
}

func (p *Pool[V]) Put(x V) {
	p.pool.Put(x)
}

func (p *Pool[V]) Get() V {
	return p.pool.Get().(V)
}
