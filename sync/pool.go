package sync

import "sync"

// Pool 是对 sync.Pool 的泛型封装。
type Pool[V any] struct {
	pool sync.Pool
	once sync.Once
	New  func() V
}

func (p *Pool[V]) init() {
	if p.New != nil {
		p.pool.New = func() any {
			return p.New()
		}
	}
}

// Put 将对象放回池中。
func (p *Pool[V]) Put(x V) {
	p.pool.Put(x)
}

// Get 从池中取出一个对象。
func (p *Pool[V]) Get() V {
	p.once.Do(p.init)

	if v, ok := p.pool.Get().(V); ok {
		return v
	}

	return *new(V)
}
