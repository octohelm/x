package sync_test

import (
	stdsync "sync"
	"testing"

	xsync "github.com/octohelm/x/sync"
	. "github.com/octohelm/x/testing/v2"
)

func TestMap(t *testing.T) {
	t.Run("零值读删与清空", func(t *testing.T) {
		var m xsync.Map[string, int]

		v, ok := m.Load("missing")
		deleted, deletedOK := m.LoadAndDelete("missing")

		Then(t, "零值 Map 对缺失键应返回类型零值和 false",
			Expect(v, Equal(0)),
			Expect(ok, Equal(false)),
			Expect(deleted, Equal(0)),
			Expect(deletedOK, Equal(false)),
			ExpectMust(func() error {
				m.Delete("missing")
				m.Clear()
				return nil
			}),
		)
	})

	t.Run("写入替换与删除语义", func(t *testing.T) {
		var m xsync.Map[string, int]

		m.Store("a", 1)
		firstLoadOrStore, loaded := m.LoadOrStore("a", 99)
		inserted, insertedLoaded := m.LoadOrStore("b", 2)
		swapped, swappedOK := m.Swap("a", 3)
		casMiss := m.CompareAndSwap("a", 1, 4)
		casHit := m.CompareAndSwap("a", 3, 4)
		cadMiss := m.CompareAndDelete("a", 3)
		cadHit := m.CompareAndDelete("a", 4)
		afterDelete, afterDeleteOK := m.Load("a")
		bValue, bOK := m.Load("b")

		Then(t, "Map 应保持标准替换与删除行为",
			Expect(firstLoadOrStore, Equal(1)),
			Expect(loaded, Equal(true)),
			Expect(inserted, Equal(2)),
			Expect(insertedLoaded, Equal(false)),
			Expect(swapped, Equal(1)),
			Expect(swappedOK, Equal(true)),
			Expect(casMiss, Equal(false)),
			Expect(casHit, Equal(true)),
			Expect(cadMiss, Equal(false)),
			Expect(cadHit, Equal(true)),
			Expect(afterDelete, Equal(0)),
			Expect(afterDeleteOK, Equal(false)),
			Expect(bValue, Equal(2)),
			Expect(bOK, Equal(true)),
		)
	})

	t.Run("遍历应支持提前停止和完整收集", func(t *testing.T) {
		var m xsync.Map[string, int]

		m.Store("a", 1)
		m.Store("b", 2)
		m.Store("c", 3)

		partial := map[string]int{}
		m.Range(func(k string, v int) bool {
			partial[k] = v
			return len(partial) < 2
		})

		full := map[string]int{}
		m.Range(func(k string, v int) bool {
			full[k] = v
			return true
		})

		Then(t, "Range 应允许提前停止且完整遍历保留所有键值",
			Expect(len(partial), Equal(2)),
			Expect(full, Equal(map[string]int{
				"a": 1,
				"b": 2,
				"c": 3,
			})),
		)
	})

	t.Run("并发读写应保持已写入值可见", func(t *testing.T) {
		var m xsync.Map[int, int]

		const workers = 32

		var wg stdsync.WaitGroup
		wg.Add(workers)

		start := make(chan struct{})
		for i := range workers {
			go func(i int) {
				defer wg.Done()
				<-start
				m.Store(i, i*i)
			}(i)
		}

		close(start)
		wg.Wait()

		Then(t, "并发写入后每个键都应能读到对应值",
			ExpectMust(func() error {
				for i := range workers {
					v, ok := m.Load(i)
					if !ok || v != i*i {
						return &loadMismatchError{Key: i, Value: v, OK: ok}
					}
				}
				return nil
			}),
		)
	})
}

type loadMismatchError struct {
	Key   int
	Value int
	OK    bool
}

func (e *loadMismatchError) Error() string {
	return "load result mismatch"
}
