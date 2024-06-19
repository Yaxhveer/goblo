package serverpool

import (
	"math/rand"
	"sync"

	"github.com/Yaxhveer/golbo/backend"
)

type randomServerPool struct {
	backends []backend.Backend
	mu       sync.RWMutex
}

func (s *randomServerPool) GetNextValidPeer() backend.Backend {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i:=0; i<5; i++ {
		rnd := rand.Intn(s.GetServerPoolSize())
		if s.backends[rnd].IsActive() {
			return s.backends[rnd]
		}
	}
	return nil
}

func (s *randomServerPool) GetBackends() []backend.Backend {
	return s.backends
}

func (s *randomServerPool) AddBackend(b backend.Backend) {
	s.backends = append(s.backends, b)
}

func (s *randomServerPool) GetServerPoolSize() int {
	return len(s.backends)
}
