package http

import (
	"fmt"
	"github.com/goal-web/contracts"
	"github.com/goal-web/pipeline"
	"github.com/goal-web/routing"
	"net/http"
)

type Engine struct {
	middlewares []contracts.MagicalFunc
	router      contracts.HttpRouter
	server      *http.Server
	app         contracts.Application
}

var (
	NotFoundResponse       = NewStringResponse(routing.NotFoundErr.Error(), 404)
	MethodNotAllowResponse = NewStringResponse(routing.MethodNotAllowErr.Error(), 405)
)

func wrapperResponse(result any) contracts.HttpResponse {
	switch v := result.(type) {
	case contracts.HttpResponse:
		return v
	case string:
		return NewStringResponse(v)
	default:
		return NewJsonResponse(v)
	}
}

func (e *Engine) ServeHTTP(writer http.ResponseWriter, r *http.Request) {
	route, params, routeErr := e.router.Route(r.Method, r.URL)
	request := NewRequest(r, params)

	if routeErr != nil {
		switch routeErr {
		case routing.MethodNotAllowErr:
			e.handleResponse(MethodNotAllowResponse, writer)
		case routing.NotFoundErr:
			e.handleResponse(NotFoundResponse, writer)
		}
		return
	}

	pipes := append(e.middlewares, route.Middlewares()...)

	var result any
	if len(pipes) == 0 {
		results := e.app.StaticCall(route.Handler(), request)
		if len(results) > 0 {
			result = results[0]
		}
	} else {
		result = pipeline.Static(e.app).
			SendStatic(request).
			ThroughStatic(pipes...).
			ThenStatic(route.Handler())
	}

	e.handleResponse(wrapperResponse(result), writer)
}

func (e *Engine) handleResponse(response contracts.HttpResponse, writer http.ResponseWriter) {
	writer.WriteHeader(response.Status())
	size, err := writer.Write(response.Bytes())
	if err != nil {
		fmt.Println(size, err)
	}
}

func (e *Engine) Start(address string) error {
	e.server = &http.Server{Addr: address, Handler: e}
	return e.server.ListenAndServe()
}

func (e *Engine) Close() error {
	return e.server.Close()
}

func (e *Engine) Static(prefix, directory string) {
}
