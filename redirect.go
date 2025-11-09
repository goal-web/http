package http

import (
	"github.com/goal-web/contracts"
	"net/http"
)

// RedirectResponse 重定向响应结构体
type RedirectResponse struct {
	*BaseResponse
	url string
}

// Redirect 创建一个重定向响应
// 参数:
//   - url: 重定向的目标URL
//   - status: HTTP状态码，默认为302
func Redirect(url string, status ...int) contracts.HttpResponse {
	redirectStatus := http.StatusFound // 默认为302
	if len(status) > 0 {
		redirectStatus = status[0]
	}

	headers := http.Header{}
	headers.Set("Location", url)

	return &RedirectResponse{
		BaseResponse: NewBaseResponse(redirectStatus, headers),
		url:          url,
	}
}

// Back 重定向到前一个页面
func Back() contracts.HttpResponse {
	headers := http.Header{}
	return &BackResponse{
		BaseResponse: NewBaseResponse(http.StatusFound, headers),
		errors:       make(map[string]string),
	}
}

// BackResponse 返回前一个页面的响应结构体
type BackResponse struct {
	*BaseResponse
	errors map[string]string
}

// WithErrors 添加错误信息到重定向响应
func (response *BackResponse) WithErrors(errors map[string]string) contracts.HttpResponse {
	response.errors = errors
	return response
}

// Bytes 实现HttpResponse接口
func (response *RedirectResponse) Bytes() []byte {
	return []byte("")
}

// WithHeaders 添加自定义响应头
func (response *RedirectResponse) WithHeaders(headers map[string]string) contracts.HttpResponse {
	response.BaseResponse.WithHeaders(headers)
	return response
}

// Header 添加单个响应头
func (response *RedirectResponse) Header(key, value string) contracts.HttpResponse {
	response.BaseResponse.Header(key, value)
	return response
}

// Cookie 添加Cookie到响应
func (response *RedirectResponse) Cookie(name, value string, maxAge int, args ...interface{}) contracts.HttpResponse {
	response.BaseResponse.Cookie(name, value, maxAge, args...)
	return response
}

// WithoutCookie 删除Cookie
func (response *RedirectResponse) WithoutCookie(name string, args ...interface{}) contracts.HttpResponse {
	response.BaseResponse.WithoutCookie(name, args...)
	return response
}

// Bytes 实现HttpResponse接口
func (response *BackResponse) Bytes() []byte {
	return []byte("")
}

// WithHeaders 添加自定义响应头
func (response *BackResponse) WithHeaders(headers map[string]string) contracts.HttpResponse {
	response.BaseResponse.WithHeaders(headers)
	return response
}

// Header 添加单个响应头
func (response *BackResponse) Header(key, value string) contracts.HttpResponse {
	response.BaseResponse.Header(key, value)
	return response
}

// Cookie 添加Cookie到响应
func (response *BackResponse) Cookie(name, value string, maxAge int, args ...interface{}) contracts.HttpResponse {
	response.BaseResponse.Cookie(name, value, maxAge, args...)
	return response
}

// WithoutCookie 删除Cookie
func (response *BackResponse) WithoutCookie(name string, args ...interface{}) contracts.HttpResponse {
	response.BaseResponse.WithoutCookie(name, args...)
	return response
}