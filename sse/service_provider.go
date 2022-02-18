package sse

import (
	"github.com/goal-web/contracts"
	"sync"
)

type ServiceProvider struct {
}

func (s ServiceProvider) Register(application contracts.Application) {
	application.Singleton("sse", func() contracts.Sse {
		return &Sse{
			fdMutex:     sync.Mutex{},
			connMutex:   sync.Mutex{},
			connections: map[uint64]contracts.SseConnection{},
			count:       0,
		}
	})
}

func (s ServiceProvider) Start() error {
	return nil
}

func (s ServiceProvider) Stop() {
}
