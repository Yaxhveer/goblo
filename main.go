package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/Yaxhveer/golbo/backend"
	"github.com/Yaxhveer/golbo/health"
	"github.com/Yaxhveer/golbo/serverpool"
	"github.com/Yaxhveer/golbo/utils"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	var configFileName string
	flag.StringVar(&configFileName, "config", "config.yaml", "config file name")

	config, err := utils.GetLBConfig(configFileName, logger)
	if err != nil {
		logger.Fatal(err.Error())
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	backends := make([]backend.Backend, 0)
	for _, u := range config.Backends {
		endpoint, err := url.Parse(u)
		if err != nil {
			logger.Fatal(err.Error(), zap.String("URL", u))
		}
		rp := httputil.NewSingleHostReverseProxy(endpoint)
		backendServer := backend.NewBackend(endpoint, rp)
		backends = append(backends, backendServer)
	}

	serverPool, err := serverpool.NewServerPool(utils.GetLBStrategy(config.Strategy), backends)
	if err != nil {
		logger.Fatal(err.Error())
	}
	loadBalancer := NewLoadBalancer(serverPool, logger)

	server := http.Server{
		Addr:    fmt.Sprintf(":%d", config.Port),
		Handler: http.HandlerFunc(loadBalancer.Serve),	
	}

	go health.LauchHealthCheck(context.Background(), serverPool, logger)

	go func() {
		<-ctx.Done()
		shutdownCtx, cancle := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancle()
		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Fatal(err)
		}
	}()

	logger.Info("Golbo started", zap.Int("port", config.Port))
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		logger.Fatal("ListenAndServe() error", zap.Error(err))
	}
}
