package websocket

import (
	"github.com/goal-web/contracts"
)

func Default(handler func(frame contracts.WebSocketFrame)) any {
	return New(&DefaultController{Handler: handler})
}

type DefaultController struct {
	Handler func(frame contracts.WebSocketFrame)
}

func (d *DefaultController) OnConnect(request contracts.HttpRequest, fd uint64) error {
	return nil
}

func (d *DefaultController) OnMessage(frame contracts.WebSocketFrame) {
	d.Handler(frame)
}

func (d *DefaultController) OnClose(fd uint64) {
}
