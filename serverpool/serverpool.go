package serverpool

import (
	"fmt"

	"github.com/Yaxhveer/golbo/backend"
	"github.com/Yaxhveer/golbo/utils"
)

type ServerPool interface {
	GetBackends() []backend.Backend
	GetNextValidPeer() backend.Backend
	AddBackend(backend.Backend)
	GetServerPoolSize() int
}

func NewServerPool(strategy utils.LBStrategy, backends []backend.Backend) (ServerPool, error) {
	switch strategy {
	case utils.RoundRobin:
		return &roundRobinServerPool{
			backends: backends,
			current:  0,
		}, nil
	case utils.LeastConnected:
		return &lcServerPool{
			backends: backends,
		}, nil
	case utils.Random:
		return &randomServerPool{
			backends: backends,
		}, nil
	default:
		return nil, fmt.Errorf("Invalid strategy")
	}
}
