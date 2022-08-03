package ferrite

import (
	"sync"
	"sync/atomic"
)

// seal is a synchronization primitive that allows synchronized calls to
// functions until the seal is closed.
type seal struct {
	done uint32 // atomic
	m    sync.Mutex
}

// Do calls fn() while holding an exclusive lock only if the seal has not been
// closed.
//
// It panics if the seal has already been closed.
func (s *seal) Do(fn func()) {
	if atomic.LoadUint32(&s.done) == 0 {
		s.m.Lock()
		defer s.m.Unlock()

		if s.done == 0 {
			fn()
			return
		}
	}

	panic("cannot modify spec after value has been used or validated")
}

// Close closes the seal.
//
// It calls fn() while holding an exclusive lock before the seal is closed.
func (s *seal) Close(fn func()) {
	if atomic.LoadUint32(&s.done) == 0 {
		s.m.Lock()
		defer s.m.Unlock()

		if s.done == 0 {
			defer atomic.StoreUint32(&s.done, 1)
			fn()
		}
	}
}
