package http

import (
	"net/http"
	"time"
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

// WithHeaders 添加多个响应头
func (base *BaseResponse) WithHeaders(headers map[string]string) *BaseResponse {
	for key, value := range headers {
		base.headers.Set(key, value)
	}
	return base
}

// Header 添加单个响应头
func (base *BaseResponse) Header(name, value string) *BaseResponse {
	base.headers.Set(name, value)
	return base
}

// Cookie 添加Cookie到响应
func (base *BaseResponse) Cookie(name, value string, maxAge int, args ...interface{}) *BaseResponse {
	cookie := &http.Cookie{
		Name:     name,
		Value:    value,
		MaxAge:   maxAge * 60, // 转换为秒
		Path:     "/",
		HttpOnly: true,
	}

	// 处理可选参数
	if len(args) > 0 {
		if path, ok := args[0].(string); ok {
			cookie.Path = path
		}
	}
	if len(args) > 1 {
		if domain, ok := args[1].(string); ok {
			cookie.Domain = domain
		}
	}
	if len(args) > 2 {
		if secure, ok := args[2].(bool); ok {
			cookie.Secure = secure
		}
	}
	if len(args) > 3 {
		if httpOnly, ok := args[3].(bool); ok {
			cookie.HttpOnly = httpOnly
		}
	}

	base.cookies = append(base.cookies, cookie)
	return base
}

// WithoutCookie 删除Cookie（通过设置过期时间为过去）
func (base *BaseResponse) WithoutCookie(name string, args ...interface{}) *BaseResponse {
	cookie := &http.Cookie{
		Name:    name,
		Value:   "",
		Expires: time.Unix(0, 0),
		MaxAge:  -1,
		Path:    "/",
	}

	// 处理可选参数
	if len(args) > 0 {
		if path, ok := args[0].(string); ok {
			cookie.Path = path
		}
	}
	if len(args) > 1 {
		if domain, ok := args[1].(string); ok {
			cookie.Domain = domain
		}
	}

	base.cookies = append(base.cookies, cookie)
	return base
}

// GetCookies 获取所有Cookie
func (base *BaseResponse) GetCookies() []*http.Cookie {
	return base.cookies
}
