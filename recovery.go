package http

import (
	"github.com/goal-web/contracts"
	"github.com/goal-web/supports/exceptions"
)

func recovery(request contracts.HttpRequest, next contracts.Pipe, handler contracts.ExceptionHandler) (result any) {
	defer func() {
		if panicValue := recover(); panicValue != nil {
			result = handler.Handle(exceptions.WrapException(panicValue))
		}
	}()
	result = next(request)
	return
}
