package http

import (
	"github.com/goal-web/contracts"
)

type StringResponse struct {
	contents string
	*BaseResponse
}

func NewStringResponse(str string, code ...int) contracts.HttpResponse {
	status := 200
	if len(code) > 0 {
		status = code[0]
	}
	return &StringResponse{contents: str, BaseResponse: &BaseResponse{
		status:  status,
		headers: map[string][]string{},
	}}
}

func (s *StringResponse) Bytes() []byte {
	return []byte(s.contents)
}
