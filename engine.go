package http

import (
	"errors"
	"fmt"
	"github.com/goal-web/container"
	"github.com/goal-web/contracts"
	"github.com/goal-web/pipeline"
	"github.com/goal-web/routing"
	"github.com/goal-web/supports/utils"
	"github.com/valyala/fasthttp"
	"net/url"
	"path/filepath"
	"strings"
)

type Engine struct {
	middlewares []contracts.MagicalFunc
	router      contracts.HttpRouter
	app         contracts.Application
	server      *fasthttp.Server

	staticDirectories map[string]string
}

func NewEngine(app contracts.Application, router contracts.HttpRouter, middlewares []contracts.MagicalFunc) contracts.HttpEngine {
	return &Engine{
		router:      router,
		app:         app,
		middlewares: middlewares,

		staticDirectories: make(map[string]string),
	}
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

func (e *Engine) tryStaticFiles(uri *url.URL, ctx *fasthttp.RequestCtx) bool {
	for prefix, directory := range e.staticDirectories {
		if strings.HasPrefix(uri.Path, prefix) {
			tryFile := filepath.Join(directory, strings.TrimPrefix(uri.Path, prefix))
			if utils.FileExists(tryFile) {
				fasthttp.ServeFile(ctx, tryFile)
				return true
			}
		}
	}
	return false
}

func (e *Engine) HandleFastHTTP(ctx *fasthttp.RequestCtx) {
	uri, _ := url.Parse(string(ctx.URI().FullURI()))
	if e.tryStaticFiles(uri, ctx) {
		return
	}

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

	pipes := append([]contracts.MagicalFunc{container.NewMagicalFunc(recovery)}, e.middlewares...)
	pipes = append(pipes, route.Middlewares()...)

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

	if result != nil {
		e.handleResponse(wrapperResponse(result), ctx)
	}
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
	e.server = &fasthttp.Server{Handler: e.HandleFastHTTP}

	return e.server.ListenAndServe(address)
}

func (e *Engine) Close() error {
	if e.server != nil {
		return e.server.Shutdown()
	}
	return nil
}

func (e *Engine) Static(prefix, directory string) {
	if !strings.HasPrefix(directory, "/") {
		directory = filepath.Join(e.app.Get("pwd").(string), directory)
	}
	e.staticDirectories[prefix] = directory
}
