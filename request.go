package http

import (
	"github.com/goal-web/contracts"
	"github.com/goal-web/supports"
	"github.com/valyala/fasthttp"
	"mime/multipart"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

var (
	defaultMemory int64 = 32 << 20 // 32 MB =
)

type Request struct {
	supports.BaseFields
	params  contracts.RouteParams
	query   *fasthttp.Args
	context map[string]any
	Request *fasthttp.RequestCtx
	fields  contracts.Fields
	lock    sync.RWMutex
}

func NewRequest(req *fasthttp.RequestCtx, params contracts.RouteParams) contracts.HttpRequest {
	var request = &Request{
		BaseFields: supports.BaseFields{},
		Request:    req,
		params:     params,
	}
	request.BaseFields.FieldsProvider = request
	request.BaseFields.OptionalGetter = request.Optional

	return request
}

func (req *Request) IsTLS() bool {
	return req.Request.IsTLS()
}

func (req *Request) GetHeader(key string) string {
	return string(req.Request.Request.Header.Peek(key))
}

func (req *Request) SetHeader(key, value string) {
	req.Request.Request.Header.Set(key, value)
}

func (req *Request) Scheme() string {
	// Can't use `r.Request.URL.Scheme`
	// See: https://groups.google.com/forum/#!topic/golang-nuts/pMUkBlQBDF0
	if req.IsTLS() {
		return "https"
	}
	if scheme := req.GetHeader("X-Forwarded-Proto"); scheme != "" {
		return scheme
	}
	if scheme := req.GetHeader("X-Forwarded-Protocol"); scheme != "" {
		return scheme
	}
	if ssl := req.GetHeader("X-Forwarded-Ssl"); ssl == "on" {
		return "https"
	}
	if scheme := req.GetHeader("X-Url-Scheme"); scheme != "" {
		return scheme
	}
	return "http"
}

func (req *Request) RealIP() string {
	// Fall back to legacy behavior
	if ip := req.GetHeader("X-Forwarded-For"); ip != "" {
		i := strings.IndexAny(ip, ",")
		if i > 0 {
			xffip := strings.TrimSpace(ip[:i])
			xffip = strings.TrimPrefix(xffip, "[")
			xffip = strings.TrimSuffix(xffip, "]")
			return xffip
		}
		return ip
	}
	if ip := req.GetHeader("X-Real-Ip"); ip != "" {
		ip = strings.TrimPrefix(ip, "[")
		ip = strings.TrimSuffix(ip, "]")
		return ip
	}
	ra, _, _ := net.SplitHostPort(req.Request.RemoteAddr().String())
	return ra
}

func (req *Request) Path() string {
	return string(req.Request.Request.URI().Path())
}

func (req *Request) Param(name string) string {
	return req.params[name]
}

func (req *Request) QueryParam(name string) string {
	if req.query == nil {
		req.query = req.Request.QueryArgs()
	}
	return string(req.query.Peek(name))
}

func (req *Request) QueryParams() url.Values {
	if req.query == nil {
		req.query = req.Request.QueryArgs()
	}
	var values = url.Values{}
	req.query.VisitAll(func(key, value []byte) {
		values.Set(string(key), string(value))
	})
	return values
}

func (req *Request) QueryString() string {
	return string(req.Request.URI().QueryString())
}

func (req *Request) FormValue(name string) string {
	return string(req.Request.FormValue(name))
}

func (req *Request) FormParams() (contracts.Fields, error) {
	form, err := req.Request.MultipartForm()
	if err != nil {
		return nil, err
	}
	var values = contracts.Fields{}
	for key, value := range form.Value {
		if len(value) > 1 {
			values[key] = value
		} else {
			values[key] = value[0]
		}
	}
	for key, file := range form.File {
		if len(file) > 1 {
			values[key] = file
		} else {
			values[key] = file[0]
		}
	}
	return values, nil
}

func (req *Request) FormFile(name string) (*multipart.FileHeader, error) {
	f, err := req.Request.FormFile(name)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (req *Request) MultipartForm() (*multipart.Form, error) {
	return req.Request.MultipartForm()
}

func (req *Request) Cookie(name string) (string, error) {
	return string(req.Request.Request.Header.Cookie(name)), nil
}

func (req *Request) SetCookie(cookie *http.Cookie) {
	req.Request.Request.Header.SetCookie(cookie.Name, cookie.Value)
}

func (req *Request) Cookies() []*http.Cookie {
	var cookies = make([]*http.Cookie, 0)
	req.Request.Request.Header.VisitAllCookie(func(key, value []byte) {
		cookies = append(cookies, &http.Cookie{
			Name:  string(key),
			Value: string(value),
		})
	})
	return cookies
}

func (req *Request) Set(key string, val interface{}) {
	req.lock.Lock()
	defer req.lock.Unlock()

	// todo
}

func (req *Request) Get(key string) any {
	return req.Optional(key, nil)
}

func (req *Request) Optional(key string, defaultValue any) (value any) {
	//if value = req.Context.Get(key); value != nil {
	//	return value
	//}
	//if value = req.Context.Param(key); value != nil && value != "" {
	//	return value
	//}
	//if req.Context.QueryParams().Has(key) {
	//	return req.Context.QueryParam(key)
	//}
	//form, err := req.Context.MultipartForm()
	//if err != nil {
	//	return defaultValue
	//}
	//if files, isFile := form.File[key]; isFile {
	//	if len(files) == 1 {
	//		return files[0]
	//	}
	//	return files
	//} else if values, isValue := form.Value[key]; isValue {
	//	if len(values) == 1 {
	//		return values[0]
	//	}
	//	return values
	//}

	return defaultValue
}

func (req *Request) Fields() contracts.Fields {
	if req.fields != nil {
		return req.fields
	}
	var data = make(contracts.Fields)

	for key, query := range req.QueryParams() {
		if len(query) == 1 {
			data[key] = query[0]
		} else {
			data[key] = query
		}
	}
	//for _, paramName := range req.ParamNames() {
	//	data[paramName] = req.Param(paramName)
	//}
	if form, existsForm := req.FormParams(); existsForm == nil {
		for key, values := range form {
			data[key] = values
		}
	}
	if multiForm, existsForm := req.MultipartForm(); existsForm == nil {
		for key, values := range multiForm.Value {
			if len(values) == 1 {
				data[key] = values[0]
			} else {
				data[key] = values
			}
		}
		for key, values := range multiForm.File {
			if len(values) == 1 {
				data[key] = values[0]
			} else {
				data[key] = values
			}
		}
	}

	req.fields = data

	return data
}
