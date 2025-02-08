package frontend

import (
	"loadbalancer-go/Serverpool"
	"net/http"
)

const (
	RETRY_ATTEMPTS int = 0
)

func AllowRetry(r *http.Request) bool {
	if _, ok := r.Context().Value(RETRY_ATTEMPTS).(bool); ok {
		return false
	}
	return true
}

type LoadBalancer interface {
	Server(http.ResponseWriter, *http.Request)
}

type loadBalancer struct {
	serverPool Serverpool.ServerPool
}

func (lb *loadBalancer) Server(rw http.ResponseWriter, req *http.Request) {
	if !AllowRetry(req) {
		rw.WriteHeader(http.StatusServiceUnavailable)
		return
	}
	peer := lb.serverPool.GetNextValidPeer()
	if peer == nil {
		rw.WriteHeader(http.StatusServiceUnavailable)
		return
	}
	peer.Serve(rw, req)
}

func NewLoadBalancer(serverPool Serverpool.ServerPool) LoadBalancer {
	return &loadBalancer{serverPool: serverPool}
}
