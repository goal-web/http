package http

import (
	"encoding/json"
	"github.com/goal-web/contracts"
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
	return &JsonResponse{content: data, BaseResponse: &BaseResponse{
		status:  status,
		headers: map[string][]string{},
	}}
}

func (s *JsonResponse) Bytes() []byte {
	jsonBytes, _ := json.Marshal(s.content)
	return jsonBytes
}
