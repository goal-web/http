package http

import (
	"github.com/goal-web/contracts"
	"github.com/goal-web/supports"
	"github.com/goal-web/supports/logs"
	"github.com/goal-web/validation"
	"github.com/labstack/echo/v4"
	"strings"
)

type Request struct {
	supports.BaseFields
	echo.Context
	fields contracts.Fields
}

func NewRequest(ctx echo.Context) contracts.HttpRequest {
	var request = &Request{
		Context:    ctx,
		fields:     nil,
		BaseFields: supports.BaseFields{},
	}

	request.BaseFields.FieldsProvider = request
	request.BaseFields.Getter = request.get

	return request
}

func (this *Request) Get(key string) (value interface{}) {
	if value = this.Context.Get(key); value != nil && value != "" {
		return value
	}
	if value = this.Context.QueryParam(key); value != nil && value != "" {
		return value
	}
	if value = this.Context.FormValue(key); value != nil && value != "" {
		return value
	}
	if value = this.Context.Param(key); value != nil && value != "" {
		return value
	}
	if file, err := this.Context.FormFile(key); err == nil && file != nil {
		return file
	}
	return
}

func (this *Request) get(key string) (value interface{}) {
	if value = this.Get(key); value != nil && value != "" {
		return value
	}
	return this.Fields()[key]
}

func (this *Request) Validate(v interface{}) error {
	if err := this.Bind(v); err != nil {
		return err
	}

	return validation.Struct(v)
}

func (this *Request) Fields() contracts.Fields {
	if this.fields != nil {
		return this.fields
	}
	var data = make(contracts.Fields)
	if strings.Contains(this.Request().Header.Get("Content-Type"), "json") {
		var bindErr = this.Context.Bind(&data)
		if bindErr != nil {
			logs.WithError(bindErr).Debug("http.Request.Fields: bind fields failed")
		}
	}

	for key, query := range this.QueryParams() {
		if len(query) == 1 {
			data[key] = query[0]
		} else {
			data[key] = query
		}
	}
	for _, paramName := range this.ParamNames() {
		data[paramName] = this.Param(paramName)
	}
	if form, existsForm := this.FormParams(); existsForm == nil {
		for key, values := range form {
			if len(values) == 1 {
				data[key] = values[0]
			} else {
				data[key] = values
			}
		}
	}
	if multiForm, existsForm := this.MultipartForm(); existsForm == nil {
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

	this.fields = data

	return data
}
