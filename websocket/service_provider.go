package websocket

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
	application.Singleton("websocket", func(config contracts.Config) contracts.WebSocket {
		var wsConfig = config.Get("websocket").(Config)

		upgrader = wsConfig.Upgrader

		return &WebSocket{
			connMutex:   sync.Mutex{},
			fdMutex:     sync.Mutex{},
			connections: map[uint64]contracts.WebSocketConnection{},
			count:       0,
		}
	})
}

func (provider ServiceProvider) Start() error {
	return nil
}

func (provider ServiceProvider) Stop() {
}
