package serverpool

import (
	"sync"

	"github.com/Yaxhveer/golbo/backend"
)

type roundRobinServerPool struct {
	backends []backend.Backend
	mu       sync.RWMutex
	current  int
}

func (s *roundRobinServerPool) rotate() backend.Backend {
	s.mu.Lock()
	defer s.mu.Unlock()
	backend := s.backends[s.current]
	s.current = (s.current + 1) % s.GetServerPoolSize()
	return backend
}

func (s *roundRobinServerPool) GetNextValidPeer() backend.Backend {
	for i := 0; i < s.GetServerPoolSize(); i++ {
		nextPeer := s.rotate()
		if nextPeer.IsActive() {
			return nextPeer
		}
	}
	return nil
}

func (s *roundRobinServerPool) GetBackends() []backend.Backend {
	return s.backends
}

func (s *roundRobinServerPool) AddBackend(b backend.Backend) {
	s.backends = append(s.backends, b)
}

func (s *roundRobinServerPool) GetServerPoolSize() int {
	return len(s.backends)
}
