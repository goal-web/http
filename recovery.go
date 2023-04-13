package http

import (
	"github.com/goal-web/contracts"
	"github.com/goal-web/supports/exceptions"
)

func (router *Router) recovery(request *Request, next contracts.Pipe) (result any) {
	defer func() {
		if panicValue := recover(); panicValue != nil {
			if res := router.errHandler(panicValue, request); res != nil { // 异常处理器返回的响应优先
				HandleResponse(res, request)
			} else {
				HandleResponse(panicValue, request) // 如果异常处理器没有定义响应，则直接响应 panic 的值
			}
			result = nil
		}
	}()
	return next(request)
}

func (router *Router) errHandler(err any, request contracts.HttpRequest) (result any) {
	var httpException Exception
	switch rawErr := err.(type) {
	case Exception:
		httpException = rawErr
	default:
		httpException = Exception{
			Exception: exceptions.WithRecover(err),
			Request:   request,
		}
	}

	// 调用容器内的异常处理器
	return router.app.StaticCall(exceptionHandler, httpException)[0]
}
