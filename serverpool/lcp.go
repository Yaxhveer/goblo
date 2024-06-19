package serverpool

import (
	"sync"

	"github.com/Yaxhveer/golbo/backend"
)

type lcServerPool struct {
	backends []backend.Backend
	mu       sync.RWMutex
}

func (s *lcServerPool) GetNextValidPeer() backend.Backend {
	var leastConnectedPeer backend.Backend
	for _, b := range s.backends {
		if b.IsActive() {
			leastConnectedPeer = b
			break
		}
	}

	for _, b := range s.backends {
		if !b.IsActive() {
			continue
		}
		if leastConnectedPeer.GetActiveConnections() > b.GetActiveConnections() {
			leastConnectedPeer = b
		}
	}
	return leastConnectedPeer
}

func (s *lcServerPool) AddBackend(b backend.Backend) {
	s.backends = append(s.backends, b)
}

func (s *lcServerPool) GetServerPoolSize() int {
	return len(s.backends)
}

func (s *lcServerPool) GetBackends() []backend.Backend {
	return s.backends
}
