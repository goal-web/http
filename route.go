package http

import "github.com/goal-web/contracts"

type route struct {
	method      []string
	path        string
	middlewares []contracts.MagicalFunc
	handler     contracts.MagicalFunc
}

func (route *route) Middlewares() []contracts.MagicalFunc {
	return route.middlewares
}

func (route *route) Method() []string {
	return route.method
}

func (route *route) Path() string {
	return route.path
}

func (route *route) Handler() contracts.MagicalFunc {
	return route.handler
}
