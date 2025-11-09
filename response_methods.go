package http

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"

	"github.com/goal-web/contracts"
)

// Response 创建一个基础响应
// 参数:
//   - content: 响应内容
//   - status: HTTP状态码（可选，默认200）
func Response(content string, status ...int) contracts.HttpResponse {
	responseStatus := 200
	if len(status) > 0 {
		responseStatus = status[0]
	}

	headers := http.Header{}
	headers.Set("Content-Type", "text/plain; charset=utf-8")

	return &GeneralResponse{
		BaseResponse: NewBaseResponse(responseStatus, headers),
		content:      content,
	}
}

// GeneralResponse 通用响应结构体
type GeneralResponse struct {
	*BaseResponse
	content string
}

func (g *GeneralResponse) Bytes() []byte {
	return []byte(g.content)
}

// Cookie 添加Cookie到响应
func (g *GeneralResponse) Cookie(name, value string, maxAge int, args ...interface{}) contracts.HttpResponse {
	g.BaseResponse.Cookie(name, value, maxAge, args...)
	return g
}

// WithoutCookie 删除Cookie
func (g *GeneralResponse) WithoutCookie(name string, args ...interface{}) contracts.HttpResponse {
	g.BaseResponse.WithoutCookie(name, args...)
	return g
}

// Header 添加单个响应头
func (g *GeneralResponse) Header(key, value string) contracts.HttpResponse {
	g.BaseResponse.Header(key, value)
	return g
}

// WithHeaders 添加自定义响应头
func (g *GeneralResponse) WithHeaders(headers map[string]string) contracts.HttpResponse {
	g.BaseResponse.WithHeaders(headers)
	return g
}

// ViewResponse 视图响应结构体
type ViewResponse struct {
	*BaseResponse
	templateName string
	data         interface{}
}

// View 创建一个视图响应
// 参数:
//   - templateName: 模板名称
//   - data: 传递给模板的数据
func View(templateName string, data ...interface{}) contracts.HttpResponse {
	headers := http.Header{}
	headers.Set("Content-Type", "text/html; charset=utf-8")

	var templateData interface{}
	if len(data) > 0 {
		templateData = data[0]
	}

	return &ViewResponse{
		BaseResponse: NewBaseResponse(200, headers),
		templateName: templateName,
		data:         templateData,
	}
}

func (v *ViewResponse) Bytes() []byte {
	// 这里应该调用视图引擎来渲染模板
	// 由于没有直接访问视图引擎，这里返回一个占位符
	return []byte(fmt.Sprintf("<!-- View: %s -->", v.templateName))
}

// WithHeaders 添加自定义响应头
func (v *ViewResponse) WithHeaders(headers map[string]string) contracts.HttpResponse {
	v.BaseResponse.WithHeaders(headers)
	return v
}

// Header 添加单个响应头
func (v *ViewResponse) Header(key, value string) contracts.HttpResponse {
	v.BaseResponse.Header(key, value)
	return v
}

// Cookie 添加Cookie到响应
func (v *ViewResponse) Cookie(name, value string, maxAge int, args ...interface{}) contracts.HttpResponse {
	v.BaseResponse.Cookie(name, value, maxAge, args...)
	return v
}

// WithoutCookie 删除Cookie
func (v *ViewResponse) WithoutCookie(name string, args ...interface{}) contracts.HttpResponse {
	v.BaseResponse.WithoutCookie(name, args...)
	return v
}

// DownloadResponse 文件下载响应结构体
type DownloadResponse struct {
	*BaseResponse
	filePath string
	fileName string
}

// Download 创建一个文件下载响应
// 参数:
//   - filePath: 文件路径
//   - fileName: 下载时显示的文件名（可选）
func Download(filePath string, fileName ...string) contracts.HttpResponse {
	headers := http.Header{}

	// 确定下载文件名
	downloadName := filepath.Base(filePath)
	if len(fileName) > 0 && fileName[0] != "" {
		downloadName = fileName[0]
	}

	// 设置下载相关的响应头
	headers.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, downloadName))
	headers.Set("Content-Type", "application/octet-stream")

	return &DownloadResponse{
		BaseResponse: NewBaseResponse(200, headers),
		filePath:     filePath,
		fileName:     downloadName,
	}
}

func (d *DownloadResponse) Bytes() []byte {
	// 读取文件内容
	content, err := os.ReadFile(d.filePath)
	if err != nil {
		// 如果文件读取失败，返回错误信息
		return []byte(fmt.Sprintf("Error reading file: %s", err.Error()))
	}
	return content
}

// WithHeaders 添加自定义响应头
func (d *DownloadResponse) WithHeaders(headers map[string]string) contracts.HttpResponse {
	d.BaseResponse.WithHeaders(headers)
	return d
}

// Header 添加单个响应头
func (d *DownloadResponse) Header(key, value string) contracts.HttpResponse {
	d.BaseResponse.Header(key, value)
	return d
}

// Cookie 添加Cookie到响应
func (d *DownloadResponse) Cookie(name, value string, maxAge int, args ...interface{}) contracts.HttpResponse {
	d.BaseResponse.Cookie(name, value, maxAge, args...)
	return d
}

// WithoutCookie 删除Cookie
func (d *DownloadResponse) WithoutCookie(name string, args ...interface{}) contracts.HttpResponse {
	d.BaseResponse.WithoutCookie(name, args...)
	return d
}

// FileResponse 文件响应结构体（用于在浏览器中显示文件）
type FileResponse struct {
	*BaseResponse
	filePath string
}

// File 创建一个文件响应（在浏览器中显示文件）
// 参数:
//   - filePath: 文件路径
func File(filePath string) contracts.HttpResponse {
	headers := http.Header{}

	// 根据文件扩展名设置Content-Type
	contentType := mime.TypeByExtension(filepath.Ext(filePath))
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	headers.Set("Content-Type", contentType)

	// 设置inline显示
	headers.Set("Content-Disposition", "inline")

	return &FileResponse{
		BaseResponse: NewBaseResponse(200, headers),
		filePath:     filePath,
	}
}

func (f *FileResponse) Bytes() []byte {
	// 读取文件内容
	content, err := os.ReadFile(f.filePath)
	if err != nil {
		// 如果文件读取失败，返回错误信息
		return []byte(fmt.Sprintf("Error reading file: %s", err.Error()))
	}
	return content
}

// WithHeaders 添加自定义响应头
func (f *FileResponse) WithHeaders(headers map[string]string) contracts.HttpResponse {
	f.BaseResponse.WithHeaders(headers)
	return f
}

// Header 添加单个响应头
func (f *FileResponse) Header(key, value string) contracts.HttpResponse {
	f.BaseResponse.Header(key, value)
	return f
}

// Cookie 添加Cookie到响应
func (f *FileResponse) Cookie(name, value string, maxAge int, args ...interface{}) contracts.HttpResponse {
	f.BaseResponse.Cookie(name, value, maxAge, args...)
	return f
}

// WithoutCookie 删除Cookie
func (f *FileResponse) WithoutCookie(name string, args ...interface{}) contracts.HttpResponse {
	f.BaseResponse.WithoutCookie(name, args...)
	return f
}

// NoContentResponse 无内容响应结构体
type NoContentResponse struct {
	*BaseResponse
}

// NoContent 创建一个无内容响应（204状态码）
func NoContent() contracts.HttpResponse {
	return &NoContentResponse{
		BaseResponse: NewBaseResponse(204, http.Header{}),
	}
}

func (n *NoContentResponse) Bytes() []byte {
	return []byte{}
}

// WithHeaders 添加自定义响应头
func (n *NoContentResponse) WithHeaders(headers map[string]string) contracts.HttpResponse {
	n.BaseResponse.WithHeaders(headers)
	return n
}

// Header 添加单个响应头
func (n *NoContentResponse) Header(key, value string) contracts.HttpResponse {
	n.BaseResponse.Header(key, value)
	return n
}

// Cookie 添加Cookie到响应
func (n *NoContentResponse) Cookie(name, value string, maxAge int, args ...interface{}) contracts.HttpResponse {
	n.BaseResponse.Cookie(name, value, maxAge, args...)
	return n
}

// WithoutCookie 删除Cookie
func (n *NoContentResponse) WithoutCookie(name string, args ...interface{}) contracts.HttpResponse {
	n.BaseResponse.WithoutCookie(name, args...)
	return n
}

// JsonpResponse JSONP响应结构体
type JsonpResponse struct {
	*JsonResponse
	callback string
}

// Jsonp 创建一个JSONP响应
// 参数:
//   - callback: 回调函数名
//   - data: 要序列化的数据
//   - code: HTTP状态码（可选，默认200）
func Jsonp(callback string, data interface{}, code ...int) contracts.HttpResponse {
	status := 200
	if len(code) > 0 {
		status = code[0]
	}

	headers := http.Header{}
	headers.Set("Content-Type", "application/javascript; charset=utf-8")

	return &JsonpResponse{
		JsonResponse: &JsonResponse{
			content:      data,
			BaseResponse: NewBaseResponse(status, headers),
		},
		callback: callback,
	}
}

func (j *JsonpResponse) Bytes() []byte {
	return []byte(fmt.Sprintf("%s(%s);", j.callback, string(j.JsonResponse.Bytes())))
}

// WithHeaders 添加自定义响应头
func (j *JsonpResponse) WithHeaders(headers map[string]string) contracts.HttpResponse {
	j.JsonResponse.WithHeaders(headers)
	return j
}

// Header 添加单个响应头
func (j *JsonpResponse) Header(key, value string) contracts.HttpResponse {
	j.JsonResponse.Header(key, value)
	return j
}

// Cookie 添加Cookie到响应
func (j *JsonpResponse) Cookie(name, value string, maxAge int, args ...interface{}) contracts.HttpResponse {
	j.JsonResponse.Cookie(name, value, maxAge, args...)
	return j
}

// WithoutCookie 删除Cookie
func (j *JsonpResponse) WithoutCookie(name string, args ...interface{}) contracts.HttpResponse {
	j.JsonResponse.WithoutCookie(name, args...)
	return j
}

// XmlResponse XML响应结构体
type XmlResponse struct {
	*BaseResponse
	content interface{}
}

// Xml 创建一个XML响应
// 参数:
//   - data: 要序列化为XML的数据
//   - code: HTTP状态码（可选，默认200）
func Xml(data interface{}, code ...int) contracts.HttpResponse {
	status := 200
	if len(code) > 0 {
		status = code[0]
	}

	headers := http.Header{}
	headers.Set("Content-Type", "application/xml; charset=utf-8")

	return &XmlResponse{
		content:      data,
		BaseResponse: NewBaseResponse(status, headers),
	}
}

func (x *XmlResponse) Bytes() []byte {
	bytes, err := xml.Marshal(x.content)
	if err != nil {
		return []byte(fmt.Sprintf("<error>XML marshal error: %s</error>", err.Error()))
	}
	return bytes
}

// WithHeaders 添加自定义响应头
func (x *XmlResponse) WithHeaders(headers map[string]string) contracts.HttpResponse {
	x.BaseResponse.WithHeaders(headers)
	return x
}

// Header 添加单个响应头
func (x *XmlResponse) Header(key, value string) contracts.HttpResponse {
	x.BaseResponse.Header(key, value)
	return x
}

// Cookie 添加Cookie到响应
func (x *XmlResponse) Cookie(name, value string, maxAge int, args ...interface{}) contracts.HttpResponse {
	x.BaseResponse.Cookie(name, value, maxAge, args...)
	return x
}

// WithoutCookie 删除Cookie
func (x *XmlResponse) WithoutCookie(name string, args ...interface{}) contracts.HttpResponse {
	x.BaseResponse.WithoutCookie(name, args...)
	return x
}

// StreamResponse 流响应结构体
type StreamResponse struct {
	*BaseResponse
	reader io.Reader
}

// Stream 创建一个流响应
// 参数:
//   - reader: 数据流读取器
//   - contentType: 内容类型（可选，默认application/octet-stream）
func Stream(reader io.Reader, contentType ...string) contracts.HttpResponse {
	headers := http.Header{}

	ct := "application/octet-stream"
	if len(contentType) > 0 && contentType[0] != "" {
		ct = contentType[0]
	}
	headers.Set("Content-Type", ct)

	return &StreamResponse{
		BaseResponse: NewBaseResponse(200, headers),
		reader:       reader,
	}
}

func (s *StreamResponse) Bytes() []byte {
	bytes, err := io.ReadAll(s.reader)
	if err != nil {
		return []byte(fmt.Sprintf("Stream read error: %s", err.Error()))
	}
	return bytes
}

// WithHeaders 添加自定义响应头
func (s *StreamResponse) WithHeaders(headers map[string]string) contracts.HttpResponse {
	s.BaseResponse.WithHeaders(headers)
	return s
}

// Header 添加单个响应头
func (s *StreamResponse) Header(key, value string) contracts.HttpResponse {
	s.BaseResponse.Header(key, value)
	return s
}

// Cookie 添加Cookie到响应
func (s *StreamResponse) Cookie(name, value string, maxAge int, args ...interface{}) contracts.HttpResponse {
	s.BaseResponse.Cookie(name, value, maxAge, args...)
	return s
}

// WithoutCookie 删除Cookie
func (s *StreamResponse) WithoutCookie(name string, args ...interface{}) contracts.HttpResponse {
	s.BaseResponse.WithoutCookie(name, args...)
	return s
}

// StatusResponse 状态响应结构体（仅返回状态码）
type StatusResponse struct {
	*BaseResponse
}

// Status 创建一个仅包含状态码的响应
// 参数:
//   - code: HTTP状态码
func Status(code int) contracts.HttpResponse {
	return &StatusResponse{
		BaseResponse: NewBaseResponse(code, http.Header{}),
	}
}

func (s *StatusResponse) Bytes() []byte {
	return []byte{}
}

// WithHeaders 添加自定义响应头
func (s *StatusResponse) WithHeaders(headers map[string]string) contracts.HttpResponse {
	s.BaseResponse.WithHeaders(headers)
	return s
}

// Header 添加单个响应头
func (s *StatusResponse) Header(key, value string) contracts.HttpResponse {
	s.BaseResponse.Header(key, value)
	return s
}

// Cookie 添加Cookie到响应
func (s *StatusResponse) Cookie(name, value string, maxAge int, args ...interface{}) contracts.HttpResponse {
	s.BaseResponse.Cookie(name, value, maxAge, args...)
	return s
}

// WithoutCookie 删除Cookie
func (s *StatusResponse) WithoutCookie(name string, args ...interface{}) contracts.HttpResponse {
	s.BaseResponse.WithoutCookie(name, args...)
	return s
}

// BinaryResponse 二进制响应结构体
type BinaryResponse struct {
	*BaseResponse
	content []byte
}

// Binary 创建一个二进制响应
// 参数:
//   - data: 二进制数据
//   - contentType: 内容类型（可选，默认application/octet-stream）
func Binary(data []byte, contentType ...string) contracts.HttpResponse {
	headers := http.Header{}

	ct := "application/octet-stream"
	if len(contentType) > 0 && contentType[0] != "" {
		ct = contentType[0]
	}
	headers.Set("Content-Type", ct)

	return &BinaryResponse{
		BaseResponse: NewBaseResponse(200, headers),
		content:      data,
	}
}

func (b *BinaryResponse) Bytes() []byte {
	return b.content
}

// WithHeaders 添加自定义响应头
func (b *BinaryResponse) WithHeaders(headers map[string]string) contracts.HttpResponse {
	b.BaseResponse.WithHeaders(headers)
	return b
}

// Header 添加单个响应头
func (b *BinaryResponse) Header(key, value string) contracts.HttpResponse {
	b.BaseResponse.Header(key, value)
	return b
}

// Cookie 添加Cookie到响应
func (b *BinaryResponse) Cookie(name, value string, maxAge int, args ...interface{}) contracts.HttpResponse {
	b.BaseResponse.Cookie(name, value, maxAge, args...)
	return b
}

// WithoutCookie 删除Cookie
func (b *BinaryResponse) WithoutCookie(name string, args ...interface{}) contracts.HttpResponse {
	b.BaseResponse.WithoutCookie(name, args...)
	return b
}

// PlainResponse 纯文本响应结构体
type PlainResponse struct {
	*BaseResponse
	content string
}

// Plain 创建一个纯文本响应
// 参数:
//   - text: 文本内容
//   - code: HTTP状态码（可选，默认200）
func Plain(text string, code ...int) contracts.HttpResponse {
	status := 200
	if len(code) > 0 {
		status = code[0]
	}

	headers := http.Header{}
	headers.Set("Content-Type", "text/plain; charset=utf-8")

	return &PlainResponse{
		BaseResponse: NewBaseResponse(status, headers),
		content:      text,
	}
}

func (p *PlainResponse) Bytes() []byte {
	return []byte(p.content)
}

// WithHeaders 添加自定义响应头
func (p *PlainResponse) WithHeaders(headers map[string]string) contracts.HttpResponse {
	p.BaseResponse.WithHeaders(headers)
	return p
}

// Header 添加单个响应头
func (p *PlainResponse) Header(key, value string) contracts.HttpResponse {
	p.BaseResponse.Header(key, value)
	return p
}

// Cookie 添加Cookie到响应
func (p *PlainResponse) Cookie(name, value string, maxAge int, args ...interface{}) contracts.HttpResponse {
	p.BaseResponse.Cookie(name, value, maxAge, args...)
	return p
}

// WithoutCookie 删除Cookie
func (p *PlainResponse) WithoutCookie(name string, args ...interface{}) contracts.HttpResponse {
	p.BaseResponse.WithoutCookie(name, args...)
	return p
}

// Text 是Plain的别名，创建一个纯文本响应
func Text(text string, code ...int) contracts.HttpResponse {
	return Plain(text, code...)
}

// HtmlResponse HTML响应结构体
type HtmlResponse struct {
	*BaseResponse
	content string
}

// Html 创建一个HTML响应
// 参数:
//   - html: HTML内容
//   - code: HTTP状态码（可选，默认200）
func Html(html string, code ...int) contracts.HttpResponse {
	status := 200
	if len(code) > 0 {
		status = code[0]
	}

	headers := http.Header{}
	headers.Set("Content-Type", "text/html; charset=utf-8")

	return &HtmlResponse{
		BaseResponse: NewBaseResponse(status, headers),
		content:      html,
	}
}

func (h *HtmlResponse) Bytes() []byte {
	return []byte(h.content)
}

// WithHeaders 添加自定义响应头
func (h *HtmlResponse) WithHeaders(headers map[string]string) contracts.HttpResponse {
	h.BaseResponse.WithHeaders(headers)
	return h
}

// Header 添加单个响应头
func (h *HtmlResponse) Header(key, value string) contracts.HttpResponse {
	h.BaseResponse.Header(key, value)
	return h
}

// Cookie 添加Cookie到响应
func (h *HtmlResponse) Cookie(name, value string, maxAge int, args ...interface{}) contracts.HttpResponse {
	h.BaseResponse.Cookie(name, value, maxAge, args...)
	return h
}

// WithoutCookie 删除Cookie
func (h *HtmlResponse) WithoutCookie(name string, args ...interface{}) contracts.HttpResponse {
	h.BaseResponse.WithoutCookie(name, args...)
	return h
}

// ErrorResponse 错误响应结构体
type ErrorResponse struct {
	*BaseResponse
	message string
	code    int
}

// Error 创建一个错误响应
// 参数:
//   - message: 错误消息
//   - code: HTTP状态码（可选，默认500）
func Error(message string, code ...int) contracts.HttpResponse {
	status := 500
	if len(code) > 0 {
		status = code[0]
	}

	headers := http.Header{}
	headers.Set("Content-Type", "application/json; charset=utf-8")

	return &ErrorResponse{
		BaseResponse: NewBaseResponse(status, headers),
		message:      message,
		code:         status,
	}
}

func (e *ErrorResponse) Bytes() []byte {
	errorData := map[string]interface{}{
		"error":   true,
		"message": e.message,
		"code":    e.code,
	}
	bytes, _ := json.Marshal(errorData)
	return bytes
}

// WithHeaders 添加自定义响应头
func (e *ErrorResponse) WithHeaders(headers map[string]string) contracts.HttpResponse {
	e.BaseResponse.WithHeaders(headers)
	return e
}

// Header 添加单个响应头
func (e *ErrorResponse) Header(key, value string) contracts.HttpResponse {
	e.BaseResponse.Header(key, value)
	return e
}

// Cookie 添加Cookie到响应
func (e *ErrorResponse) Cookie(name, value string, maxAge int, args ...interface{}) contracts.HttpResponse {
	e.BaseResponse.Cookie(name, value, maxAge, args...)
	return e
}

// WithoutCookie 删除Cookie
func (e *ErrorResponse) WithoutCookie(name string, args ...interface{}) contracts.HttpResponse {
	e.BaseResponse.WithoutCookie(name, args...)
	return e
}

// Abort 创建一个中止响应（通常用于中间件）
// 参数:
//   - code: HTTP状态码
//   - message: 错误消息（可选）
func Abort(code int, message ...string) contracts.HttpResponse {
	msg := http.StatusText(code)
	if len(message) > 0 && message[0] != "" {
		msg = message[0]
	}

	headers := http.Header{}
	headers.Set("Content-Type", "application/json; charset=utf-8")

	return &ErrorResponse{
		BaseResponse: NewBaseResponse(code, headers),
		message:      msg,
		code:         code,
	}
}

// NotFound 创建一个404响应
// 参数:
//   - message: 错误消息（可选）
func NotFound(message ...string) contracts.HttpResponse {
	msg := "Not Found"
	if len(message) > 0 && message[0] != "" {
		msg = message[0]
	}
	return Error(msg, 404)
}

// Forbidden 创建一个403响应
// 参数:
//   - message: 错误消息（可选）
func Forbidden(message ...string) contracts.HttpResponse {
	msg := "Forbidden"
	if len(message) > 0 && message[0] != "" {
		msg = message[0]
	}
	return Error(msg, 403)
}

// Unauthorized 创建一个401响应
// 参数:
//   - message: 错误消息（可选）
func Unauthorized(message ...string) contracts.HttpResponse {
	msg := "Unauthorized"
	if len(message) > 0 && message[0] != "" {
		msg = message[0]
	}
	return Error(msg, 401)
}

// BadRequest 创建一个400响应
// 参数:
//   - message: 错误消息（可选）
func BadRequest(message ...string) contracts.HttpResponse {
	msg := "Bad Request"
	if len(message) > 0 && message[0] != "" {
		msg = message[0]
	}
	return Error(msg, 400)
}

// InternalServerError 创建一个500响应
// 参数:
//   - message: 错误消息（可选）
func InternalServerError(message ...string) contracts.HttpResponse {
	msg := "Internal Server Error"
	if len(message) > 0 && message[0] != "" {
		msg = message[0]
	}
	return Error(msg, 500)
}