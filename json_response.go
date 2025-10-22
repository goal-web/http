package http

import (
	"encoding/json"
	"github.com/goal-web/contracts"
	"net/http"
)

type JsonResponse struct {
	content any
	*BaseResponse
}

func NewJsonResponse(data any, code ...int) contracts.HttpResponse {
	status := 200
	if len(code) > 0 {
		status = code[0]
	}
	headers := http.Header{}
	headers.Set("Content-Type", "application/json; charset=utf-8")
	return &JsonResponse{
		content:      data,
		BaseResponse: NewBaseResponse(status, headers),
	}
}

func (s *JsonResponse) Bytes() []byte {
	jsonBytes, _ := json.Marshal(s.content)
	return jsonBytes
}

// WithHeaders 添加自定义响应头
func (s *JsonResponse) WithHeaders(headers map[string]string) contracts.HttpResponse {
	s.BaseResponse.WithHeaders(headers)
	return s
}

// Header 添加单个响应头
func (s *JsonResponse) Header(key, value string) contracts.HttpResponse {
	s.BaseResponse.Header(key, value)
	return s
}

// Cookie 添加Cookie到响应
func (s *JsonResponse) Cookie(name, value string, maxAge int, args ...interface{}) contracts.HttpResponse {
	s.BaseResponse.Cookie(name, value, maxAge, args...)
	return s
}

// WithoutCookie 删除Cookie
func (s *JsonResponse) WithoutCookie(name string, args ...interface{}) contracts.HttpResponse {
	s.BaseResponse.WithoutCookie(name, args...)
	return s
}

// Json 创建JSON响应的简化方法
// 参数:
//   - data: 要序列化的数据
//   - code: HTTP状态码（可选，默认200）
func Json(data interface{}, code ...int) contracts.HttpResponse {
	return NewJsonResponse(data, code...)
}
