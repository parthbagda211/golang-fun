package Serverpool

import (
	"context"
	_ "context"
	"fmt"
	_ "fmt"
	"go.uber.org/zap"
	"loadbalancer-go/Utils"
	"sync"
	_ "sync"
	"time"
	_ "time"

	_ "go.uber.org/zap"
	_ "loadbalancer-go/Utils"
	"loadbalancer-go/backend"
)

type ServerPool interface {
	GetBackend() []backend.Backend
	GetNextValidPeer() backend.Backend
	AddBackend(backend.Backend)
	GetServerPoolSize() int
}

type rounRobinServerPool struct {
	backends []backend.Backend
	mux      sync.RWMutex
	current  int
}

func (s *rounRobinServerPool) Rotate() backend.Backend {
	s.mux.Lock()
	s.current = (s.current + 1) % s.GetServerPoolSize()
	s.mux.Unlock()
	return s.backends[s.current]

}
func (s *rounRobinServerPool) GetNextValidPeer() backend.Backend {
	for i := 0; i < s.GetServerPoolSize(); i++ {
		nextPeer := s.Rotate()
		if nextPeer.IsAlive() {
			return nextPeer
		}
	}
	return nil
}

func (s *rounRobinServerPool) GetServerPoolSize() int {
	return len(s.backends)
}

func (s *rounRobinServerPool) GetBackend() []backend.Backend {
	return s.backends
}

func (s *rounRobinServerPool) AddBackend(b backend.Backend) {
	s.mux.Lock()
	s.backends = append(s.backends, b)
	s.mux.Unlock()
}

func HealthCheck(ctx context.Context, s ServerPool) {
	aliveChannel := make(chan bool, 1)
	for _, b := range s.GetBackend() {
		b := b
		requestCtx, stop := context.WithTimeout(ctx, 10*time.Second)
		defer stop()
		status := "up"
		go backend.IsBackendAlive(requestCtx, aliveChannel, b.GetUrl())

		select {
		case <-ctx.Done():
			Utils.Logger.Info("Shutting Down Health Check")
			return
		case alive := <-aliveChannel:
			b.SetAlive(alive)
			if !alive {
				status = "down"
			}
		}
		Utils.Logger.Debug("URL Status",
			zap.String("url", b.GetUrl().String()),
			zap.String("status", status),
		)
	}

}

func NewServerPool(stategy Utils.LBStrategy) (ServerPool, error) {
	switch stategy {
	case Utils.RoundRobin:
		return &rounRobinServerPool{
			backends: make([]backend.Backend, 0),
			current:  0,
		}, nil
	case Utils.LeastConnected:
		return &lcServerPool{
			backends: make([]backend.Backend, 0),
		}, nil
	default:
		return nil, fmt.Errorf("Invalid Stategy")
	}
}
