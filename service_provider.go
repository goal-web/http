package http

import (
	"fmt"
	"github.com/goal-web/contracts"
	"github.com/goal-web/routing"
	"github.com/goal-web/supports/logs"
	"github.com/goal-web/supports/utils"
	"net/http"
)

type ServiceProvider struct {
	app             contracts.Application
	engine          contracts.HttpEngine
	RouteCollectors []any
}

func NewService(routes ...any) contracts.ServiceProvider {
	return &ServiceProvider{RouteCollectors: routes}
}

func (provider *ServiceProvider) Stop() {
	provider.app.Call(func(dispatcher contracts.EventDispatcher) {
		if err := provider.engine.Close(); err != nil {
			logs.WithError(err).Info("failed to close http engine.")
		}
		dispatcher.Dispatch(&ServeClosed{})
	})
}

func (provider *ServiceProvider) Start() error {
	for _, collector := range provider.RouteCollectors {
		provider.app.Call(collector)
	}

	var err error
	provider.app.Call(func(
		router contracts.HttpRouter,
		config contracts.Config,
		events contracts.EventDispatcher,
	) {
		httpConfig := config.Get("http").(Config)

		provider.engine = &Engine{
			router:      router,
			middlewares: append(routing.ConvertToMiddlewares(httpConfig.GlobalMiddlewares...), router.Middlewares()...),
			app:         provider.app,
		}

		for prefix, directory := range httpConfig.StaticDirectories {
			provider.engine.Static(prefix, directory)
		}

		err = router.Mount()
		if err != nil {
			return
		}

		err = provider.engine.Start(
			utils.StringOr(
				httpConfig.Address,
				fmt.Sprintf("%s:%s", httpConfig.Host, utils.StringOr(httpConfig.Port, "8000")),
			),
		)
	})

	if err != nil && err != http.ErrServerClosed {
		logs.WithError(err).Error("http service failed to start")
		go func() { provider.app.Stop() }()
		return err
	}

	return nil
}

func (provider *ServiceProvider) Register(app contracts.Application) {
	provider.app = app

	app.Singleton("HttpRouter", func() contracts.HttpRouter {
		return routing.NewHttpRouter(provider.app)
	})
}
