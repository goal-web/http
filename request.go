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
	var request = &Request{Context: ctx, BaseFields: supports.BaseFields{}}
	request.BaseFields.FieldsProvider = request
	request.BaseFields.OptionalGetter = request.Optional

	return request
}

func (request *Request) Get(key string) any {
	return request.Optional(key, nil)
}

func (request *Request) Optional(key string, defaultValue any) (value any) {
	if value = request.Context.Get(key); value != nil {
		return value
	}
	if value = request.Context.Param(key); value != nil && value != "" {
		return value
	}
	if request.Context.QueryParams().Has(key) {
		return request.Context.QueryParam(key)
	}
	form, err := request.Context.MultipartForm()
	if err != nil {
		return defaultValue
	}
	if files, isFile := form.File[key]; isFile {
		if len(files) == 1 {
			return files[0]
		}
		return files
	} else if values, isValue := form.Value[key]; isValue {
		if len(values) == 1 {
			return values[0]
		}
		return values
	}

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
	if strings.Contains(request.Request().Header.Get("Content-Type"), "json") {
		var bindErr = request.Context.Bind(&data)
		if bindErr != nil {
			logs.WithError(bindErr).Debug("http.Request.Fields: bind fields failed")
		}
	}

	for key, query := range request.QueryParams() {
		if len(query) == 1 {
			data[key] = query[0]
		} else {
			data[key] = query
		}
	}
	for _, paramName := range request.ParamNames() {
		data[paramName] = request.Param(paramName)
	}
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
