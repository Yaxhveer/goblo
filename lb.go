package main

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/Yaxhveer/golbo/serverpool"
)

type LoadBalancer interface {
	Serve(http.ResponseWriter, *http.Request)
}

type loadBalancer struct {
	serverPool serverpool.ServerPool
	logger     *zap.Logger
}

func (lb *loadBalancer) Serve(w http.ResponseWriter, r *http.Request) {
	peer := lb.serverPool.GetNextValidPeer()
	if peer != nil {
		lb.logger.Info("access to endpoint",
			zap.String("url", peer.GetURL().String()),
			zap.Int("connections", peer.GetActiveConnections()),
		)
		peer.Serve(w, r)
		return
	}
	http.Error(w, "Service not available", http.StatusServiceUnavailable)
}

func NewLoadBalancer(serverPool serverpool.ServerPool, logger *zap.Logger) LoadBalancer {
	return &loadBalancer{
		serverPool: serverPool,
		logger:     logger,
	}
}
