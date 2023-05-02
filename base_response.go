package http

import "net/http"

type BaseResponse struct {
	status  int
	headers http.Header
}

func (base *BaseResponse) Status() int {
	return base.status
}

func (base *BaseResponse) Headers() http.Header {
	return base.headers
}

func (base *BaseResponse) SetStatus(status int) {
	base.status = status
}

func (base *BaseResponse) AddHeader(name string, header string) {
	base.headers[name] = append(base.headers[name], header)
}

func (base *BaseResponse) DelHeader(name string) {
	delete(base.headers, name)
}

func (base *BaseResponse) SetHeader(name string, headers []string) {
	base.headers[name] = headers
}

func (base *BaseResponse) SetHeaders(headers http.Header) {
	base.headers = headers
}
