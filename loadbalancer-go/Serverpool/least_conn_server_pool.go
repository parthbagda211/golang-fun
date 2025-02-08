package Serverpool

import (
	"loadbalancer-go/backend"
	"sync"
)

type lcServerPool struct {
	backends []backend.Backend
	mux      sync.RWMutex
}

func (s *lcServerPool) AddBacked(b backend.Backend) {
	s.backends = append(s.backends, b)
}

func (s *lcServerPool) getServerPoolSize() int {
	return len(s.backends)
}

func (s *lcServerPool) GetBackends() []backend.Backend {
	return s.backends
}
func (s *lcServerPool) GetNextValidPeer() backend.Backend {
	var leastConnectedPeer backend.Backend
	for _, b := range s.backends {
		if b.IsAlive() {
			leastConnectedPeer = b
			break
		}
	}

	for _, b := range s.backends {
		if !b.IsAlive() {
			continue
		}
		if leastConnectedPeer.GetActiveConnections() > b.GetActiveConnections() {
			leastConnectedPeer = b
		}
	}
	return leastConnectedPeer
}
