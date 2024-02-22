package sse

import (
	"github.com/goal-web/contracts"
	"sync"
)

type Factory struct {
	drivers map[string]contracts.Sse
	mutex   sync.RWMutex
}

func NewFactory() contracts.SseFactory {
	return &Factory{
		drivers: map[string]contracts.Sse{},
		mutex:   sync.RWMutex{},
	}
}

func (factory *Factory) Sse(key string) contracts.Sse {
	factory.mutex.RLock()
	defer factory.mutex.RUnlock()
	driver, _ := factory.drivers[key]
	return driver
}

func (factory *Factory) Register(key string, sse contracts.Sse) {
	factory.mutex.Lock()
	defer factory.mutex.Unlock()
	factory.drivers[key] = sse
}
