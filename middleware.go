package http

import (
	"fmt"
	"github.com/goal-web/container"
	"github.com/goal-web/contracts"
	"github.com/goal-web/supports/exceptions"
	"sync"
)

type Middleware struct {
	container contracts.Container

	middlewares map[string]contracts.MagicalFunc

	mutex sync.Mutex
}

func NewMiddleware(container contracts.Container) contracts.Middleware {
	return &Middleware{
		middlewares: make(map[string]contracts.MagicalFunc),
		container:   container,
		mutex:       sync.Mutex{},
	}
}

func (m *Middleware) Register(name string, handler any) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if _, exists := m.middlewares[name]; exists {
		panic(MiddlewareDuplicateException{
			Exception: exceptions.New("duplicate middleware: " + name),
			Name:      name,
			Handler:   handler,
		})
	}
	m.middlewares[name] = container.NewMagicalFunc(handler)
}

type MiddlewareNotFoundException struct {
	contracts.Exception

	Name      string
	Arguments []any
}

type MiddlewareDuplicateException struct {
	contracts.Exception

	Name    string
	Handler any
}

func (m *Middleware) Call(name string, params ...any) any {
	if middleware, exists := m.middlewares[name]; exists {
		return m.container.StaticCall(middleware, params...)[0]
	}

	panic(MiddlewareNotFoundException{
		Exception: exceptions.New(fmt.Sprintf("middleware %s not found", name)),
		Name:      name,
		Arguments: params,
	})
}
