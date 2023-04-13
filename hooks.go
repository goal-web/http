package http

import "github.com/goal-web/contracts"

type RequestBefore struct {
	request contracts.HttpRequest
}

func (event *RequestBefore) Event() string {
	return "REQUEST_BEFORE"
}

func (event *RequestBefore) Sync() bool {
	return true
}
func (event *RequestBefore) Request() contracts.HttpRequest {
	return event.request
}

type RequestAfter struct {
	request contracts.HttpRequest
}

func (event *RequestAfter) Event() string {
	return "REQUEST_AFTER"
}

func (event *RequestAfter) Request() contracts.HttpRequest {
	return event.request
}

type ResponseBefore struct {
	request contracts.HttpRequest
}

func (events *ResponseBefore) Event() string {
	return "RESPONSE_BEFORE"
}

func (events *ResponseBefore) Request() contracts.HttpRequest {
	return events.request
}

func (events *ResponseBefore) Sync() bool {
	return true
}
