package backend

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
)

type Backend interface {
	SetActive(bool)
	IsActive() bool
	GetURL() *url.URL
	GetActiveConnections() int
	Serve(http.ResponseWriter, *http.Request)
}

type backend struct {
	url          *url.URL
	active       bool
	mu           sync.RWMutex
	connections  int
	reverseProxy *httputil.ReverseProxy
}

func (b *backend) GetActiveConnections() int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	connections := b.connections
	return connections
}

func (b *backend) SetActive(active bool) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.active = active
}

func (b *backend) IsActive() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	active := b.active
	return active
}

func (b *backend) GetURL() *url.URL {
	return b.url
}

func (b *backend) Serve(w http.ResponseWriter, r *http.Request) {
	defer func() {
		b.mu.Lock()
		b.connections--
		b.mu.Unlock()
	}()

	b.mu.Lock()
	b.connections++
	b.mu.Unlock()
	b.reverseProxy.ServeHTTP(w, r)
}

func NewBackend(u *url.URL, rp *httputil.ReverseProxy) Backend {
	return &backend{
		url:          u,
		active:       false,
		reverseProxy: rp,
	}
}
