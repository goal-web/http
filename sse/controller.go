package sse

import (
	"github.com/goal-web/contracts"
)

func Default() interface{} {
	return New(&DefaultController{})
}

type DefaultController struct {
}

func (d *DefaultController) OnConnect(request contracts.HttpRequest, fd uint64) error {
	return nil
}

func (d *DefaultController) OnClose(fd uint64) {
}
