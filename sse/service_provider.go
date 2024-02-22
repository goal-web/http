package sse

import (
	"github.com/goal-web/contracts"
)

type ServiceProvider struct {
}

func NewService() contracts.ServiceProvider {
	return &ServiceProvider{}
}

func (provider ServiceProvider) Register(application contracts.Application) {
	application.Singleton("sse", func() contracts.Sse {
		return NewSse()
	})
	application.Singleton("sse.factory", func() contracts.SseFactory {
		return NewFactory()
	})
}

func (provider ServiceProvider) Start() error {
	return nil
}

func (provider ServiceProvider) Stop() {
}
