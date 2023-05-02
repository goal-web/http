package http

import (
	"github.com/goal-web/contracts"
)

func recovery(request contracts.HttpRequest, next contracts.Pipe) (result any) {
	defer func() {
		if panicValue := recover(); panicValue != nil {
			// todo
		}
	}()
	return next(request)
}
