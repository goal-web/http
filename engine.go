package http

import (
	"errors"
	"fmt"
	"github.com/goal-web/contracts"
	"github.com/goal-web/pipeline"
	"github.com/goal-web/routing"
	"github.com/valyala/fasthttp"
	"net/url"
	"strings"
)

type Engine struct {
	middlewares []contracts.MagicalFunc
	router      contracts.HttpRouter
	app         contracts.Application
	server      *fasthttp.Server
}

func (e *Engine) Request() contracts.HttpRequest {
	//TODO implement me
	panic("implement me")
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

func (e *Engine) HandleFastHTTP(ctx *fasthttp.RequestCtx) {
	uri, _ := url.Parse(string(ctx.URI().FullURI()))
	route, params, routeErr := e.router.Route(string(ctx.Method()), uri)
	request := NewRequest(ctx, params)

	if routeErr != nil {
		switch {
		case errors.Is(routeErr, routing.MethodNotAllowErr):
			e.handleResponse(MethodNotAllowResponse, ctx)
		case errors.Is(routeErr, routing.NotFoundErr):
			e.handleResponse(NotFoundResponse, ctx)
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

	e.handleResponse(wrapperResponse(result), ctx)
}

func (e *Engine) handleResponse(response contracts.HttpResponse, ctx *fasthttp.RequestCtx) {
	for key, value := range response.Headers() {
		if len(value) > 1 {
			ctx.Response.Header.Set(key, strings.Join(value, ";"))
		} else {
			ctx.Response.Header.Set(key, value[0])
		}
	}
	ctx.SetStatusCode(response.Status())
	size, err := ctx.Write(response.Bytes())
	if err != nil {
		fmt.Println(size, err)
	}
}

func (e *Engine) Start(address string) error {
	e.server = &fasthttp.Server{
		Handler: e.HandleFastHTTP,
	}

	return e.server.ListenAndServe(address)
}

func (e *Engine) Close() error {
	return e.server.Shutdown()
}

func (e *Engine) Static(prefix, directory string) {
}
