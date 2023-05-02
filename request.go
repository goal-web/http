package http

import (
	"github.com/goal-web/contracts"
	"github.com/goal-web/supports"
	"github.com/goal-web/validation"
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
	query   url.Values
	context map[string]any
	request *http.Request
	http.ResponseWriter
	fields contracts.Fields
	lock   sync.RWMutex
}

func NewRequest(req *http.Request, params contracts.RouteParams) contracts.HttpRequest {
	var request = &Request{
		BaseFields: supports.BaseFields{},
		request:    req,
		params:     params,
	}
	request.BaseFields.FieldsProvider = request
	request.BaseFields.OptionalGetter = request.Optional

	return request
}

func (request *Request) Request() *http.Request {
	return request.request
}

func (request *Request) IsTLS() bool {
	return request.request.TLS != nil
}

func (request *Request) IsWebSocket() bool {
	upgrade := request.request.Header.Get("Upgrade")
	return strings.EqualFold(upgrade, "websocket")
}

func (request *Request) Scheme() string {
	// Can't use `r.Request.URL.Scheme`
	// See: https://groups.google.com/forum/#!topic/golang-nuts/pMUkBlQBDF0
	if request.IsTLS() {
		return "https"
	}
	if scheme := request.request.Header.Get("X-Forwarded-Proto"); scheme != "" {
		return scheme
	}
	if scheme := request.request.Header.Get("X-Forwarded-Protocol"); scheme != "" {
		return scheme
	}
	if ssl := request.request.Header.Get("X-Forwarded-Ssl"); ssl == "on" {
		return "https"
	}
	if scheme := request.request.Header.Get("X-Url-Scheme"); scheme != "" {
		return scheme
	}
	return "http"
}

func (request *Request) RealIP() string {
	// Fall back to legacy behavior
	if ip := request.request.Header.Get("X-Forwarded-For"); ip != "" {
		i := strings.IndexAny(ip, ",")
		if i > 0 {
			xffip := strings.TrimSpace(ip[:i])
			xffip = strings.TrimPrefix(xffip, "[")
			xffip = strings.TrimSuffix(xffip, "]")
			return xffip
		}
		return ip
	}
	if ip := request.request.Header.Get("X-Real-Ip"); ip != "" {
		ip = strings.TrimPrefix(ip, "[")
		ip = strings.TrimSuffix(ip, "]")
		return ip
	}
	ra, _, _ := net.SplitHostPort(request.request.RemoteAddr)
	return ra
}

func (request *Request) Path() string {
	return request.request.URL.Path
}

func (request *Request) Param(name string) string {
	return request.params[name]
}

func (request *Request) QueryParam(name string) string {
	if request.query == nil {
		request.query = request.request.URL.Query()
	}
	return request.query.Get(name)
}

func (request *Request) QueryParams() url.Values {
	if request.query == nil {
		request.query = request.request.URL.Query()
	}
	return request.request.URL.Query()
}

func (request *Request) QueryString() string {
	return request.request.URL.RawQuery
}

func (request *Request) FormValue(name string) string {
	return request.request.FormValue(name)
}

func (request *Request) FormParams() (url.Values, error) {
	//if strings.HasPrefix(request.request.Header.Get(HeaderContentType), MIMEMultipartForm) {
	//	if err := request.request.ParseMultipartForm(defaultMemory); err != nil {
	//		return nil, err
	//	}
	//} else {
	//	if err := request.request.ParseForm(); err != nil {
	//		return nil, err
	//	}
	//}
	return request.request.Form, nil
}

func (request *Request) FormFile(name string) (*multipart.FileHeader, error) {
	f, fh, err := request.request.FormFile(name)
	if err != nil {
		return nil, err
	}
	f.Close()
	return fh, nil
}

func (request *Request) MultipartForm() (*multipart.Form, error) {
	err := request.request.ParseMultipartForm(defaultMemory)
	return request.request.MultipartForm, err
}

func (request *Request) Cookie(name string) (*http.Cookie, error) {
	return request.request.Cookie(name)
}

func (request *Request) SetCookie(cookie *http.Cookie) {
	// todo
}

func (request *Request) Cookies() []*http.Cookie {
	return request.request.Cookies()
}

func (request *Request) Set(key string, val interface{}) {
	request.lock.Lock()
	defer request.lock.Unlock()

	// todo
}

func (request *Request) Bind(i interface{}) error {
	return nil
}

func (request *Request) Get(key string) any {
	return request.Optional(key, nil)
}

func (request *Request) Optional(key string, defaultValue any) (value any) {
	//if value = request.Context.Get(key); value != nil {
	//	return value
	//}
	//if value = request.Context.Param(key); value != nil && value != "" {
	//	return value
	//}
	//if request.Context.QueryParams().Has(key) {
	//	return request.Context.QueryParam(key)
	//}
	//form, err := request.Context.MultipartForm()
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

func (request *Request) Validate(v any) error {
	if err := request.Bind(v); err != nil {
		return err
	}

	return validation.Struct(v)
}

func (request *Request) Fields() contracts.Fields {
	if request.fields != nil {
		return request.fields
	}
	var data = make(contracts.Fields)

	for key, query := range request.QueryParams() {
		if len(query) == 1 {
			data[key] = query[0]
		} else {
			data[key] = query
		}
	}
	//for _, paramName := range request.ParamNames() {
	//	data[paramName] = request.Param(paramName)
	//}
	if form, existsForm := request.FormParams(); existsForm == nil {
		for key, values := range form {
			if len(values) == 1 {
				data[key] = values[0]
			} else {
				data[key] = values
			}
		}
	}
	if multiForm, existsForm := request.MultipartForm(); existsForm == nil {
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

	request.fields = data

	return data
}
