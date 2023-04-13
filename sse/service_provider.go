package sse

import (
	"github.com/goal-web/contracts"
	"sync"
)

type ServiceProvider struct {
}

func NewService() contracts.ServiceProvider {
	return &ServiceProvider{}
}

func (provider ServiceProvider) Register(application contracts.Application) {
	application.Singleton("sse", func() contracts.Sse {
		return &Sse{
			fdMutex:     sync.Mutex{},
			connMutex:   sync.Mutex{},
			connections: map[uint64]contracts.SseConnection{},
			count:       0,
		}
	})
}

func (provider ServiceProvider) Start() error {
	return nil
}

func (provider ServiceProvider) Stop() {
}
