// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"reflect"
	"sync"
	"sync/atomic"
	"testing"
)

type bench struct {
	setup func(*testing.B, mapInterface)
	perG  func(b *testing.B, pb *testing.PB, i int, m mapInterface)
}

func benchMap(b *testing.B, bench bench) {
	for _, m := range [...]mapInterface{&DeepCopyMap{}, &RWMutexMap{}, &sync.Map{}} {
		b.Run(fmt.Sprintf("%T", m), func(b *testing.B) {
			m = reflect.New(reflect.TypeOf(m).Elem()).Interface().(mapInterface)
			if bench.setup != nil {
				bench.setup(b, m)
			}

			b.ResetTimer()

			var i int64
			b.RunParallel(func(pb *testing.PB) {
				id := int(atomic.AddInt64(&i, 1) - 1)
				bench.perG(b, pb, id*b.N, m)
			})
		})
	}
}

func BenchmarkLoadMostlyHits(b *testing.B) {
	const hits, misses = 1023, 1

	benchMap(b, bench{
		setup: func(_ *testing.B, m mapInterface) {
			for i := 0; i < hits; i++ {
				m.LoadOrStore(i, i)
			}
			// Prime the map to get it into a steady state.
			for i := 0; i < hits*2; i++ {
				m.Load(i % hits)
			}
		},

		perG: func(b *testing.B, pb *testing.PB, i int, m mapInterface) {
			for ; pb.Next(); i++ {
				m.Load(i % (hits + misses))
			}
		},
	})

	// syncmap code:
	b.Run(fmt.Sprintf("%T", &IntMap{}), func(b *testing.B) {
		m := &IntMap{}
		// setup:
		for i := 0; i < hits; i++ {
			m.LoadOrStore(i, i)
		}
		for i := 0; i < hits*2; i++ {
			m.Load(i % hits)
		}

		// reset:
		b.ResetTimer()

		// perG:
		perG := func(b *testing.B, pb *testing.PB, i int, m *IntMap) {
			for ; pb.Next(); i++ {
				m.Load(i % (hits + misses))
			}
		}
		var i int64
		b.RunParallel(func(pb *testing.PB) {
			id := int(atomic.AddInt64(&i, 1) - 1)
			perG(b, pb, id*b.N, m)
		})
	})
}

func BenchmarkLoadMostlyMisses(b *testing.B) {
	const hits, misses = 1, 1023

	benchMap(b, bench{
		setup: func(_ *testing.B, m mapInterface) {
			for i := 0; i < hits; i++ {
				m.LoadOrStore(i, i)
			}
			// Prime the map to get it into a steady state.
			for i := 0; i < hits*2; i++ {
				m.Load(i % hits)
			}
		},

		perG: func(b *testing.B, pb *testing.PB, i int, m mapInterface) {
			for ; pb.Next(); i++ {
				m.Load(i % (hits + misses))
			}
		},
	})

	// syncmap code:
	b.Run(fmt.Sprintf("%T", &IntMap{}), func(b *testing.B) {
		m := &IntMap{}
		// setup:
		for i := 0; i < hits; i++ {
			m.LoadOrStore(i, i)
		}
		for i := 0; i < hits*2; i++ {
			m.Load(i % hits)
		}

		// reset:
		b.ResetTimer()

		// perG:
		perG := func(b *testing.B, pb *testing.PB, i int, m *IntMap) {
			for ; pb.Next(); i++ {
				m.Load(i % (hits + misses))
			}
		}
		var i int64
		b.RunParallel(func(pb *testing.PB) {
			id := int(atomic.AddInt64(&i, 1) - 1)
			perG(b, pb, id*b.N, m)
		})
	})
}

func BenchmarkLoadOrStoreBalanced(b *testing.B) {
	const hits, misses = 128, 128

	benchMap(b, bench{
		setup: func(b *testing.B, m mapInterface) {
			if _, ok := m.(*DeepCopyMap); ok {
				b.Skip("DeepCopyMap has quadratic running time.")
			}
			for i := 0; i < hits; i++ {
				m.LoadOrStore(i, i)
			}
			// Prime the map to get it into a steady state.
			for i := 0; i < hits*2; i++ {
				m.Load(i % hits)
			}
		},

		perG: func(b *testing.B, pb *testing.PB, i int, m mapInterface) {
			for ; pb.Next(); i++ {
				j := i % (hits + misses)
				if j < hits {
					if _, ok := m.LoadOrStore(j, i); !ok {
						b.Fatalf("unexpected miss for %v", j)
					}
				} else {
					if v, loaded := m.LoadOrStore(i, i); loaded {
						b.Fatalf("failed to store %v: existing value %v", i, v)
					}
				}
			}
		},
	})

	// syncmap code:
	b.Run(fmt.Sprintf("%T", &IntMap{}), func(b *testing.B) {
		m := &IntMap{}
		// setup:
		for i := 0; i < hits; i++ {
			m.LoadOrStore(i, i)
		}
		for i := 0; i < hits*2; i++ {
			m.Load(i % hits)
		}

		// reset:
		b.ResetTimer()

		// perG:
		perG := func(b *testing.B, pb *testing.PB, i int, m *IntMap) {
			for ; pb.Next(); i++ {
				j := i % (hits + misses)
				if j < hits {
					if _, ok := m.LoadOrStore(j, i); !ok {
						b.Fatalf("unexpected miss for %v", j)
					}
				} else {
					if v, loaded := m.LoadOrStore(i, i); loaded {
						b.Fatalf("failed to store %v: existing value %v", i, v)
					}
				}
			}
		}

		var i int64
		b.RunParallel(func(pb *testing.PB) {
			id := int(atomic.AddInt64(&i, 1) - 1)
			perG(b, pb, id*b.N, m)
		})
	})
}

func BenchmarkLoadOrStoreUnique(b *testing.B) {
	benchMap(b, bench{
		setup: func(b *testing.B, m mapInterface) {
			if _, ok := m.(*DeepCopyMap); ok {
				b.Skip("DeepCopyMap has quadratic running time.")
			}
		},

		perG: func(b *testing.B, pb *testing.PB, i int, m mapInterface) {
			for ; pb.Next(); i++ {
				m.LoadOrStore(i, i)
			}
		},
	})

	// syncmap code:
	b.Run(fmt.Sprintf("%T", &IntMap{}), func(b *testing.B) {
		m := &IntMap{}
		// setup:

		// reset:
		b.ResetTimer()

		// perG:
		perG := func(b *testing.B, pb *testing.PB, i int, m *IntMap) {
			for ; pb.Next(); i++ {
				m.LoadOrStore(i, i)
			}
		}
		var i int64
		b.RunParallel(func(pb *testing.PB) {
			id := int(atomic.AddInt64(&i, 1) - 1)
			perG(b, pb, id*b.N, m)
		})
	})
}

func BenchmarkLoadOrStoreCollision(b *testing.B) {
	benchMap(b, bench{
		setup: func(_ *testing.B, m mapInterface) {
			m.LoadOrStore(0, 0)
		},

		perG: func(b *testing.B, pb *testing.PB, i int, m mapInterface) {
			for ; pb.Next(); i++ {
				m.LoadOrStore(0, 0)
			}
		},
	})

	// syncmap code:
	b.Run(fmt.Sprintf("%T", &IntMap{}), func(b *testing.B) {
		m := &IntMap{}
		// setup:

		// reset:
		b.ResetTimer()

		// perG:
		perG := func(b *testing.B, pb *testing.PB, i int, m *IntMap) {
			for ; pb.Next(); i++ {
				m.LoadOrStore(0, 0)
			}
		}
		var i int64
		b.RunParallel(func(pb *testing.PB) {
			id := int(atomic.AddInt64(&i, 1) - 1)
			perG(b, pb, id*b.N, m)
		})
	})
}

func BenchmarkRange(b *testing.B) {
	const mapSize = 1 << 10

	benchMap(b, bench{
		setup: func(_ *testing.B, m mapInterface) {
			for i := 0; i < mapSize; i++ {
				m.Store(i, i)
			}
		},

		perG: func(b *testing.B, pb *testing.PB, i int, m mapInterface) {
			for ; pb.Next(); i++ {
				m.Range(func(_, _ interface{}) bool { return true })
			}
		},
	})

	// syncmap code:
	b.Run(fmt.Sprintf("%T", &IntMap{}), func(b *testing.B) {
		m := &IntMap{}
		// setup:

		// reset:
		b.ResetTimer()

		// perG:
		perG := func(b *testing.B, pb *testing.PB, i int, m *IntMap) {
			for ; pb.Next(); i++ {
				m.Range(func(_, _ int) bool { return true })
			}
		}
		var i int64
		b.RunParallel(func(pb *testing.PB) {
			id := int(atomic.AddInt64(&i, 1) - 1)
			perG(b, pb, id*b.N, m)
		})
	})
}

// BenchmarkAdversarialAlloc tests performance when we store a new value
// immediately whenever the map is promoted to clean and otherwise load a
// unique, missing key.
//
// This forces the Load calls to always acquire the map's mutex.
func BenchmarkAdversarialAlloc(b *testing.B) {
	benchMap(b, bench{
		perG: func(b *testing.B, pb *testing.PB, i int, m mapInterface) {
			var stores, loadsSinceStore int64
			for ; pb.Next(); i++ {
				m.Load(i)
				if loadsSinceStore++; loadsSinceStore > stores {
					m.LoadOrStore(i, stores)
					loadsSinceStore = 0
					stores++
				}
			}
		},
	})

	// syncmap code:
	b.Run(fmt.Sprintf("%T", &IntMap{}), func(b *testing.B) {
		m := &IntMap{}
		// setup:

		// reset:
		b.ResetTimer()

		// perG:
		perG := func(b *testing.B, pb *testing.PB, i int, m *IntMap) {
			var stores, loadsSinceStore int
			for ; pb.Next(); i++ {
				m.Load(i)
				if loadsSinceStore++; loadsSinceStore > stores {
					m.LoadOrStore(i, stores)
					loadsSinceStore = 0
					stores++
				}
			}
		}
		var i int64
		b.RunParallel(func(pb *testing.PB) {
			id := int(atomic.AddInt64(&i, 1) - 1)
			perG(b, pb, id*b.N, m)
		})
	})
}

// BenchmarkAdversarialDelete tests performance when we periodically delete
// one key and add a different one in a large map.
//
// This forces the Load calls to always acquire the map's mutex and periodically
// makes a full copy of the map despite changing only one entry.
func BenchmarkAdversarialDelete(b *testing.B) {
	const mapSize = 1 << 10

	benchMap(b, bench{
		setup: func(_ *testing.B, m mapInterface) {
			for i := 0; i < mapSize; i++ {
				m.Store(i, i)
			}
		},

		perG: func(b *testing.B, pb *testing.PB, i int, m mapInterface) {
			for ; pb.Next(); i++ {
				m.Load(i)

				if i%mapSize == 0 {
					m.Range(func(k, _ interface{}) bool {
						m.Delete(k)
						return false
					})
					m.Store(i, i)
				}
			}
		},
	})

	// syncmap code:
	b.Run(fmt.Sprintf("%T", &IntMap{}), func(b *testing.B) {
		m := &IntMap{}
		// setup:

		// reset:
		b.ResetTimer()

		// perG:
		perG := func(b *testing.B, pb *testing.PB, i int, m *IntMap) {
			for ; pb.Next(); i++ {
				m.Load(i)

				if i%mapSize == 0 {
					m.Range(func(k, _ int) bool {
						m.Delete(k)
						return false
					})
					m.Store(i, i)
				}
			}
		}
		var i int64
		b.RunParallel(func(pb *testing.PB) {
			id := int(atomic.AddInt64(&i, 1) - 1)
			perG(b, pb, id*b.N, m)
		})
	})
}
