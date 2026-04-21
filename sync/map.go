package sync

import "sync"

// Map 是对 sync.Map 的泛型封装。
type Map[K comparable, V any] struct {
	m sync.Map
}

func (m *Map[K, V]) typeAssertion(v any, result bool) (V, bool) {
	if vv, ok := v.(V); ok {
		return vv, result
	}
	return *new(V), result
}

// LoadOrStore 返回现有值，或在键不存在时写入并返回给定值。
func (m *Map[K, V]) LoadOrStore(k K, value V) (V, bool) {
	return m.typeAssertion(m.m.LoadOrStore(k, value))
}

// Delete 删除指定键。
func (m *Map[K, V]) Delete(k K) {
	m.m.Delete(k)
}

// Load 读取指定键对应的值。
func (m *Map[K, V]) Load(k K) (V, bool) {
	return m.typeAssertion(m.m.Load(k))
}

// LoadAndDelete 读取并删除指定键。
func (m *Map[K, V]) LoadAndDelete(k K) (V, bool) {
	return m.typeAssertion(m.m.LoadAndDelete(k))
}

// Swap 将键更新为新值，并返回旧值。
func (m *Map[K, V]) Swap(k K, value V) (V, bool) {
	return m.typeAssertion(m.m.Swap(k, value))
}

// CompareAndSwap 在旧值匹配时写入新值。
func (m *Map[K, V]) CompareAndSwap(k K, old V, new V) bool {
	return m.m.CompareAndSwap(k, old, new)
}

// CompareAndDelete 在旧值匹配时删除键。
func (m *Map[K, V]) CompareAndDelete(k K, old V) bool {
	return m.m.CompareAndDelete(k, old)
}

// Store 写入指定键值。
func (m *Map[K, V]) Store(k K, v V) {
	m.m.Store(k, v)
}

// Clear 清空整个 Map。
func (m *Map[K, V]) Clear() {
	m.m.Clear()
}

// Range 依次遍历 Map 中的所有键值对。
func (m *Map[K, V]) Range(r func(k K, v V) bool) {
	m.m.Range(func(key, value any) bool {
		return r(key.(K), value.(V))
	})
}
