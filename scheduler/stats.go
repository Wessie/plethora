package scheduler

import (
	"sync"
	"sync/atomic"
)

func NewStats() *Stats {
	return &Stats{
		events: make(map[string]*uint64),
	}
}

type Stats struct {
	mu     sync.RWMutex
	events map[string]*uint64
}

func (s *Stats) Inc(name string) {
	s.mu.RLock()
	i, ok := s.events[name]
	s.mu.RUnlock()

	if !ok {
		s.mu.Lock()
		i, ok = s.events[name]
		if !ok {
			vi := uint64(1)
			s.events[name] = &vi
		}
		s.mu.Unlock()
	}

	atomic.AddUint64(i, 1)
}

func (s *Stats) IncTask(tsk Task) {
	s.Inc(toString(tsk.Job))
	s.Inc(toString(tsk.Planner))
	s.Inc(tsk.String())
}
