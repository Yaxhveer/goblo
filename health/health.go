package health

import (
	"context"
	"net"
	"net/url"
	"time"

	"go.uber.org/zap"

	"github.com/Yaxhveer/golbo/serverpool"
)

func IsBackendActive(ctx context.Context, ActiveChannel chan bool, u *url.URL, logger *zap.Logger) {
	var d net.Dialer
	conn, err := d.DialContext(ctx, "tcp", u.Host)
	if err != nil {
		logger.Info("Site unreachable", zap.Error(err))
		ActiveChannel <- false
		return
	}
	_ = conn.Close()
	ActiveChannel <- true
}

func HealthCheck(ctx context.Context, s serverpool.ServerPool, logger *zap.Logger) {
	ActiveChannel := make(chan bool, 1)

	for _, b := range s.GetBackends() {
		b := b
		requestCtx, stop := context.WithTimeout(ctx, 10*time.Second)
		defer stop()
		status := "up"
		go IsBackendActive(requestCtx, ActiveChannel, b.GetURL(), logger)

		select {
		case <-ctx.Done():
			logger.Info("Gracefully shutting down health check")
			return
		case active := <-ActiveChannel:
			if b.IsActive() != active {
				b.SetActive(active)
			}
			if !active {
				status = "down"
			}
		}
		logger.Info(
			"URL Status",
			zap.String("URL", b.GetURL().String()),
			zap.String("status", status),
		)
	}
}

func LauchHealthCheck(ctx context.Context, sp serverpool.ServerPool, logger *zap.Logger) {
	logger.Info("Starting health check...")
	go HealthCheck(ctx, sp, logger)
	t := time.NewTicker(time.Second * 15)
	for {
		select {
		case <-t.C:
			go HealthCheck(ctx, sp, logger)
		case <-ctx.Done():
			logger.Info("Closing Health Check")
			return
		}
	}
}
