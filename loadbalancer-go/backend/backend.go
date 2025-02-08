package backend

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
)

type Backend interface {
	SetAlive(bool)
	IsAlive() bool
	GetUrl() *url.URL
	GetActiveConnections() int
	Serve(rw http.ResponseWriter, req *http.Request)
}

type backend struct {
	url         *url.URL
	alive       bool
	mux         sync.RWMutex
	connections int
	reveseProxy *httputil.ReverseProxy
}

/*
Acquires a read lock on the mux mutex to ensure thread-safe access to the connections field.
Reads the value of the connections field.
Releases the read lock on the mux mutex.
*/
func (b *backend) GetActiveConnections() int {
	b.mux.RLock()
	connections := b.connections
	b.mux.RUnlock()
	return connections
}

func (b *backend) SetAlive(alive bool) {
	b.mux.Lock()
	b.alive = alive
	b.mux.Unlock()
}

func (b *backend) IsAlive() bool {
	b.mux.RLock()
	alive := b.alive
	b.mux.RUnlock()
	return alive
}

func (b *backend) GetUrl() *url.URL {
	return b.url
}

func (b *backend) Serve(rw http.ResponseWriter, req *http.Request) {
	defer func() {
		b.mux.Lock()
		b.connections--
		b.mux.Unlock()
	}()

	b.mux.Lock()
	b.connections++
	b.mux.Unlock()
	b.reveseProxy.ServeHTTP(rw, req)
}

func NewBackend(u *url.URL, rp *httputil.ReverseProxy) Backend {
	return &backend{
		url:         u,
		alive:       true,
		reveseProxy: rp,
	}
}
