package ferrite

import (
	"sync"
	"sync/atomic"
)

type smutex struct {
	done uint32 // atomic
	m    sync.RWMutex
}

func (s *smutex) Lock() {
	if atomic.LoadUint32(&s.done) == 0 {
		s.m.Lock()

		if s.done == 0 {
			return
		}

		s.m.Unlock()
	}

	panic("cannot modify spec after value has been used or validated")
}

func (s *smutex) Unlock() {
	s.m.Unlock()
}

func (s *smutex) RLock() {
	if atomic.LoadUint32(&s.done) == 0 {
		s.m.RLock()
	}
}

func (s *smutex) RUnlock() {
	if atomic.LoadUint32(&s.done) == 0 {
		s.m.RUnlock()
	}
}

func (s *smutex) Seal(fn func()) {
	if atomic.LoadUint32(&s.done) == 0 {
		s.m.Lock()
		defer s.m.Unlock()

		if s.done == 0 {
			fn()
			atomic.StoreUint32(&s.done, 1)
		}
	}
}
