package http

import (
	"errors"
	"github.com/goal-web/container"
	"github.com/goal-web/contracts"
	"github.com/goal-web/pipeline"
	"github.com/goal-web/supports/exceptions"
	"github.com/labstack/echo/v4"
)

var (
	MiddlewareError = errors.New("middleware error") // 中间件必须有一个返回值

	// magical middlewares
	exceptionHandler = container.NewMagicalFunc(func(handler contracts.ExceptionHandler, exception Exception) any {
		return handler.Handle(exception)
	})
)

func New(container contracts.Application) contracts.Router {
	router := &Router{
		app:         container,
		events:      container.Get("events").(contracts.EventDispatcher),
		echo:        echo.New(),
		routes:      make([]contracts.Route, 0),
		groups:      make([]contracts.RouteGroup, 0),
		middlewares: make([]contracts.MagicalFunc, 0),
	}

	router.Use(router.recovery)

	return router
}

type Router struct {
	events contracts.EventDispatcher
	app    contracts.Application
	echo   *echo.Echo
	groups []contracts.RouteGroup
	routes []contracts.Route

	// 全局中间件
	middlewares []contracts.MagicalFunc
}

func (router *Router) Group(prefix string, middlewares ...any) contracts.RouteGroup {
	groupInstance := NewGroup(prefix, middlewares...)

	router.groups = append(router.groups, groupInstance)

	return groupInstance
}

func (router *Router) Close() error {
	return router.echo.Close()
}

func (router *Router) Static(path, directory string) {
	router.echo.Static(path, directory)
}

func (router *Router) Get(path string, handler any, middlewares ...any) {
	router.Add(echo.GET, path, handler, middlewares...)
}

func (router *Router) Post(path string, handler any, middlewares ...any) {
	router.Add(echo.POST, path, handler, middlewares...)
}

func (router *Router) Delete(path string, handler any, middlewares ...any) {
	router.Add(echo.DELETE, path, handler, middlewares...)
}

func (router *Router) Put(path string, handler any, middlewares ...any) {
	router.Add(echo.PUT, path, handler, middlewares...)
}

func (router *Router) Patch(path string, handler any, middlewares ...any) {
	router.Add(echo.PATCH, path, handler, middlewares...)
}

func (router *Router) Options(path string, handler any, middlewares ...any) {
	router.Add(echo.OPTIONS, path, handler, middlewares...)
}

func (router *Router) Trace(path string, handler any, middlewares ...any) {
	router.Add(echo.TRACE, path, handler, middlewares...)
}

func (router *Router) Use(middlewares ...any) {
	for _, middleware := range middlewares {
		if magicalFunc, ok := middleware.(contracts.MagicalFunc); ok {
			router.middlewares = append(router.middlewares, magicalFunc)
		} else if echoMiddleware, isEchoFunc := middleware.(echo.MiddlewareFunc); isEchoFunc {
			router.echo.Use(echoMiddleware)
		} else {
			router.middlewares = append(router.middlewares, container.NewMagicalFunc(middleware))
		}
	}
}

func (router *Router) Add(method any, path string, handler any, middlewares ...any) {
	methods := make([]string, 0)
	switch v := method.(type) {
	case string:
		methods = []string{v}
	case []string:
		methods = v
	default:
		panic(errors.New("method 只能接收 string 或者 []string"))
	}
	router.routes = append(router.routes, &route{
		method:      methods,
		path:        path,
		middlewares: convertToMiddlewares(middlewares...),
		handler:     container.NewMagicalFunc(handler),
	})
}

func (router *Router) mountGroup(group contracts.RouteGroup) {
	router.mountRoutes(group.Routes(), group.Middlewares()...)

	for _, routeGroup := range group.Groups() {
		router.mountGroup(routeGroup)
	}
}

// Start 启动 httpserver
func (router *Router) Start(address string) error {

	router.mountRoutes(router.routes)

	for _, routeGroup := range router.groups {
		router.mountGroup(routeGroup)
	}

	router.echo.HTTPErrorHandler = func(err error, context echo.Context) {
		if result := router.app.StaticCall(exceptionHandler, Exception{Exception: exceptions.WithError(err), Request: NewRequest(context)})[0]; result != nil {
			HandleResponse(result, NewRequest(context))
		}
	}
	router.echo.Debug = router.app.Debug()

	return router.echo.Start(address)
}

// mountRoutes 装配路由
func (router *Router) mountRoutes(routes []contracts.Route, middlewares ...contracts.MagicalFunc) {
	for _, routeItem := range routes {
		(func(routeInstance contracts.Route) {
			router.echo.Match(routeInstance.Method(), routeInstance.Path(), func(context echo.Context) error {
				request := NewRequest(context)
				defer func() {
					router.events.Dispatch(&RequestAfter{request})
				}()

				// 触发钩子
				router.events.Dispatch(&RequestBefore{request})

				pipes := append(router.middlewares, middlewares...)
				pipes = append(pipes, routeInstance.Middlewares()...)

				var result any
				if len(pipes) == 0 {
					results := router.app.StaticCall(routeInstance.Handler(), request)
					if len(results) > 0 {
						result = results[0]
					}
				} else {
					result = pipeline.Static(router.app).SendStatic(request).
						ThroughStatic(
							router.middlewares...,
						).
						ThroughStatic(
							append(middlewares, routeInstance.Middlewares()...)...,
						).
						ThenStatic(routeInstance.Handler())
				}

				router.events.Dispatch(&ResponseBefore{request})

				HandleResponse(result, request)

				return nil
			})
		})(routeItem)
	}
}
