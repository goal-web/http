package http

import (
	"github.com/goal-web/contracts"
	"net/http"
)

type StringResponse struct {
	contents string
	*BaseResponse
	cookies []*http.Cookie
}

func NewStringResponse(str string, code ...int) contracts.HttpResponse {
	status := 200
	if len(code) > 0 {
		status = code[0]
	}
	headers := http.Header{}
	headers.Set("Content-Type", "text/plain; charset=utf-8")
	return &StringResponse{
		contents:     str,
		BaseResponse: NewBaseResponse(status, headers),
		cookies:      make([]*http.Cookie, 0),
	}
}

func (s *StringResponse) Bytes() []byte {
	return []byte(s.contents)
}

// WithHeaders 添加自定义响应头
func (s *StringResponse) WithHeaders(headers map[string]string) contracts.HttpResponse {
	for key, value := range headers {
		s.Headers().Set(key, value)
	}
	return s
}

// Header 添加单个响应头
func (s *StringResponse) Header(key, value string) contracts.HttpResponse {
	s.Headers().Set(key, value)
	return s
}
