package http

import (
	"github.com/goal-web/contracts"
)

type Exception struct {
	contracts.Exception
	Request contracts.HttpRequest
}
