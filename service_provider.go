package http

import (
	"fmt"
	"github.com/goal-web/contracts"
	"github.com/goal-web/routing"
	"github.com/goal-web/supports/logs"
	"github.com/goal-web/supports/utils"
	"github.com/pkg/errors"
	"net/http"
)

type ServiceProvider struct {
	app             contracts.Application
	RouteCollectors []any
}

func NewService(routes ...any) contracts.ServiceProvider {
	return &ServiceProvider{RouteCollectors: routes}
}

func (provider *ServiceProvider) Stop() {
	provider.app.Call(func(dispatcher contracts.EventDispatcher, engine contracts.HttpEngine) {
		if err := engine.Close(); err != nil {
			logs.WithError(err).Info("failed to close http engine.")
		}
		dispatcher.Dispatch(&ServeClosed{})
	})
}

func (provider *ServiceProvider) Start() error {

	var err error
	provider.app.Call(func(
		router contracts.HttpRouter,
		engine contracts.HttpEngine,
		config contracts.Config,
		events contracts.EventDispatcher,
	) {
		httpConfig := config.Get("http").(Config)

		for prefix, directory := range httpConfig.StaticDirectories {
			engine.Static(prefix, directory)
		}

		err = router.Mount()
		if err != nil {
			return
		}

		router.Print()
		address := utils.StringOr(
			httpConfig.Address,
			fmt.Sprintf("%s:%s", httpConfig.Host, utils.StringOr(httpConfig.Port, "8000")),
		)
		fmt.Println("http service is started.")
		fmt.Println(address)
		err = engine.Start(address)
	})

	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		logs.WithError(err).Error("http service failed to start")
		go func() { provider.app.Stop() }()
		return err
	}

	return nil
}

func (provider *ServiceProvider) Register(app contracts.Application) {
	provider.app = app

	app.Singleton("HttpMiddleware", func() contracts.Middleware {
		return NewMiddleware(provider.app)
	})
	app.Singleton("HttpRouter", func() contracts.HttpRouter {
		return routing.NewHttpRouter(provider.app)
	})
	app.Singleton("HttpEngine", func(router contracts.HttpRouter, config contracts.Config, middleware contracts.Middleware) contracts.HttpEngine {
		httpConfig := config.Get("http").(Config)
		return NewEngine(provider.app, router, append(routing.ConvertToMiddlewares(middleware, httpConfig.GlobalMiddlewares...), router.Middlewares()...))
	})
	app.Call(func(console contracts.Console) {
		console.RegisterCommand(routing.NewRouteList)
	})

	for _, collector := range provider.RouteCollectors {
		provider.app.Call(collector)
	}
}
