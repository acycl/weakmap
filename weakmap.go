// Package weakmap provides a generic cache map implementation with [weak
// pointers](https://pkg.go.dev/weak#Pointer) as described [on the Go
// blog](https://go.dev/blog/cleanups-and-weak#weak-pointers).
package weakmap

import (
	"runtime"
	"sync"
	"weak"
)

// Map is a concurrent cache that stores values as weak pointers. Values are
// created on demand via the New function and automatically removed from the
// cache when they are no longer referenced elsewhere and garbage collected.
//
// A Map must not be copied after first use.
type Map[K comparable, V any] struct {
	// New creates a value for the given key. It is called when Load is invoked
	// for a key that is not in the cache or whose cached value has been
	// collected. If New is nil, Load always returns nil.
	New   func(K) *V
	cache sync.Map
}

// Load retrieves a value from the cache, creating it via New if necessary.
// If the cached value has been garbage collected, New is called again to
// create a fresh value. Returns nil if New is nil.
func (m *Map[K, V]) Load(key K) *V {
	if m.New == nil {
		return nil
	}

	var v *V
	for {
		iwp, loaded := m.cache.Load(key)
		if !loaded {

			// Only create a new value once.
			if v == nil {
				v = m.New(key)
			}

			wp := weak.Make(v)
			iwp, loaded = m.cache.LoadOrStore(key, wp)
			if !loaded {
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
