package http

import (
	"encoding/json"
	"fmt"
	"github.com/valyala/fasthttp"
	"reflect"
	"strconv"
	"strings"
)

// ParseRequestParams 用于解析请求的参数到指定结构体，包括表单、文件和 JSON 数据
func ParseRequestParams(ctx *fasthttp.RequestCtx, dest interface{}) error {
	// 获取目标结构体的类型和值
	destValue := reflect.ValueOf(dest).Elem()
	destType := destValue.Type()

	// 默认解析 URL 查询参数
	if args := ctx.QueryArgs(); args != nil {
		if err := parseFormValues("query", args, destValue, destType); err != nil {
			return err
		}
	}

	// 判断 Content-Type
	contentType := string(ctx.Request.Header.ContentType())

	// 根据不同的 Content-Type 类型处理不同的数据解析
	switch {
	case strings.Contains(contentType, "application/json"):
		// 解析 JSON 数据
		if err := json.Unmarshal(ctx.PostBody(), dest); err != nil {
			return fmt.Errorf("failed to parse JSON: %w", err)
		}
	case strings.Contains(contentType, "application/x-www-form-urlencoded"):
		// 解析表单数据
		args := ctx.PostArgs()
		if err := parseFormValues("form", args, destValue, destType); err != nil {
			return err
		}
	case strings.Contains(contentType, "multipart/form-data"):
		// 解析多部分表单（包括文件）
		if err := parseMultipartForm(ctx, destValue, destType); err != nil {
			return err
		}

	}

	return nil
}

// parseFormValues 解析 URL 查询参数或 application/x-www-form-urlencoded 参数
func parseFormValues(tagName string, args *fasthttp.Args, destValue reflect.Value, destType reflect.Type) error {
	for i := 0; i < destValue.NumField(); i++ {
		field := destValue.Field(i)
		fieldType := destType.Field(i)
		tagValue := fieldType.Tag.Get(tagName)

		if tagValue == "" {
			tagValue = fieldType.Name // 使用字段名作为默认参数名
		}

		paramValue := args.Peek(tagValue)
		if len(paramValue) == 0 {
			continue
		}

		if err := setValue(field, string(paramValue)); err != nil {
			return fmt.Errorf("failed to set value for field %s: %w", fieldType.Name, err)
		}
	}
	return nil
}

// parseMultipartForm 解析 multipart/form-data 表单，包括文件
func parseMultipartForm(ctx *fasthttp.RequestCtx, destValue reflect.Value, destType reflect.Type) error {
	// 解析 multipart form
	multipartForm, err := ctx.MultipartForm()
	if err != nil {
		return fmt.Errorf("failed to parse multipart/form-data: %w", err)
	}

	// 处理表单字段
	for key, values := range multipartForm.Value {
		for i := 0; i < destValue.NumField(); i++ {
			field := destValue.Field(i)
			fieldType := destType.Field(i)
			tagValue := fieldType.Tag.Get("form")

			if tagValue == "" {
				tagValue = fieldType.Name
			}

			if key == tagValue {
				if err = setValue(field, values[0]); err != nil {
					return fmt.Errorf("failed to set form value for field %s: %w", fieldType.Name, err)
				}
			}
		}
	}

	// 处理文件字段
	for key, files := range multipartForm.File {
		for i := 0; i < destValue.NumField(); i++ {
			field := destValue.Field(i)
			fieldType := destType.Field(i)

			tagValue := fieldType.Tag.Get("form")
			if tagValue == "" {
				tagValue = fieldType.Name
			}

			if key == tagValue {
				if field.Kind() == reflect.Slice || field.Kind() == reflect.Array {
					field.Set(reflect.ValueOf(files))
				} else {
					field.Set(reflect.ValueOf(files[0]))
				}
			}
		}
	}

	return nil
}

// setValue 根据字段类型设置对应的值
func setValue(field reflect.Value, value string) error {
	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		intVal, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		field.SetInt(intVal)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		uintVal, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return err
		}
		field.SetUint(uintVal)
	case reflect.Float32, reflect.Float64:
		floatVal, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}
		field.SetFloat(floatVal)
	case reflect.Bool:
		boolVal, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		field.SetBool(boolVal)
	default:
		return fmt.Errorf("unsupported kind %s", field.Kind())
	}
	return nil
}
