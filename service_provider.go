package http

import (
	"fmt"
	"github.com/goal-web/contracts"
	"github.com/goal-web/supports/logs"
	"github.com/goal-web/supports/utils"
	"net/http"
)

type ServiceProvider struct {
	app contracts.Application

	RouteCollectors []any
}

func NewService(routes ...any) contracts.ServiceProvider {
	return &ServiceProvider{RouteCollectors: routes}
}

func (provider *ServiceProvider) Stop() {
	provider.app.Call(func(dispatcher contracts.EventDispatcher, router contracts.Router) {
		if err := router.Close(); err != nil {
			logs.WithError(err).Info("Router 关闭报错")
		}
		dispatcher.Dispatch(&ServeClosed{})
	})
}

func (provider *ServiceProvider) Start() error {
	for _, collector := range provider.RouteCollectors {
		provider.app.Call(collector)
	}

	var err error
	provider.app.Call(func(router contracts.Router, config contracts.Config) {
		httpConfig := config.Get("http").(Config)
		err = router.Start(
			utils.StringOr(
				httpConfig.Address,
				fmt.Sprintf("%s:%s", httpConfig.Host, utils.StringOr(httpConfig.Port, "8000")),
			),
		)
	})

	if err != nil && err != http.ErrServerClosed {
		logs.WithError(err).Error("http 服务无法启动")
		go func() { provider.app.Stop() }()
		return err
	}

	return nil
}

func (provider *ServiceProvider) Register(app contracts.Application) {
	provider.app = app

	app.Singleton("Router", func() contracts.Router {
		return New(provider.app)
	})
}
