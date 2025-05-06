// Package weakmap provides a generic cache map implementation with [weak
// pointers](https://pkg.go.dev/weak#Pointer) as described [on the Go
// blog](https://go.dev/blog/cleanups-and-weak#weak-pointers).
package weakmap

import (
	"runtime"
	"sync"
	"weak"
)

type Map[K comparable, V any] struct {
	New   func(K) *V
	cache sync.Map
}

func (m *Map[K, V]) Load(key K) *V {
	if m.New == nil {
		return nil
	}

	var v *V
	for {
		iwp, ok := m.cache.Load(key)
		if !ok {

			// Only create a new value once.
			if v == nil {
				v = m.New(key)
			}

			wp := weak.Make(v)
			iwp, ok = m.cache.LoadOrStore(key, wp)
			if ok {
				runtime.AddCleanup(v, func(key K) {
					m.cache.CompareAndDelete(key, wp)
				}, key)
				return v
			}
		}

		if v := iwp.(weak.Pointer[V]).Value(); v != nil {
			return v
		}

		m.cache.CompareAndDelete(key, iwp)
	}
}
