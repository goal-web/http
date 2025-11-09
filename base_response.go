package http

import (
	"net/http"
)

type BaseResponse struct {
	status  int
	headers http.Header
	cookies []*http.Cookie
}

func NewBaseResponse(status int, headers http.Header) *BaseResponse {
	if headers == nil {
		headers = make(http.Header)
	}
	return &BaseResponse{
		status:  status,
		headers: headers,
		cookies: make([]*http.Cookie, 0),
	}
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

// WithHeaders 添加自定义响应头
func (base *BaseResponse) WithHeaders(headers map[string]string) *BaseResponse {
	for key, value := range headers {
		base.headers.Set(key, value)
	}
	return base
}

// Header 添加单个响应头
func (base *BaseResponse) Header(key, value string) *BaseResponse {
	base.headers.Set(key, value)
	return base
}

// Cookie 添加Cookie到响应
func (base *BaseResponse) Cookie(name, value string, maxAge int, args ...interface{}) *BaseResponse {
	cookie := &http.Cookie{
		Name:  name,
		Value: value,
		Path:  "/",
	}

	if maxAge > 0 {
		cookie.MaxAge = maxAge
	}

	// 处理其他参数
	if len(args) > 0 {
		if domain, ok := args[0].(string); ok {
			cookie.Domain = domain
		}
	}

	if len(args) > 1 {
		if path, ok := args[1].(string); ok {
			cookie.Path = path
		}
	}

	// 添加或更新Cookie
	found := false
	for i, c := range base.cookies {
		if c.Name == name {
			base.cookies[i] = cookie
			found = true
			break
		}
	}
	if !found {
		base.cookies = append(base.cookies, cookie)
	}

	return base
}

// WithoutCookie 删除Cookie
func (base *BaseResponse) WithoutCookie(name string, args ...interface{}) *BaseResponse {
	for i, cookie := range base.cookies {
		if cookie.Name == name {
			// 创建新的cookies切片，排除要删除的cookie
			base.cookies = append(base.cookies[:i], base.cookies[i+1:]...)
			break
		}
	}
	return base
}
