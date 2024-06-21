package serverpool

import (
	"math/rand"

	"github.com/Yaxhveer/golbo/backend"
)

type randomServerPool struct {
	backends []backend.Backend
}

func (s *randomServerPool) GetNextValidPeer() backend.Backend {
	for i:=0; i<3; i++ {
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
