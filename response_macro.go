package http

import (
	"github.com/goal-web/contracts"
	"sync"
)

// ResponseMacroFunc 响应宏函数类型
type ResponseMacroFunc func(args ...interface{}) contracts.HttpResponse

// 全局响应宏注册表
var (
	responseMacros = make(map[string]ResponseMacroFunc)
	macroMutex     sync.RWMutex
)

// ResponseMacro 注册一个响应宏
// 参数:
//   - name: 宏名称
//   - macro: 宏函数
func ResponseMacro(name string, macro ResponseMacroFunc) {
	macroMutex.Lock()
	defer macroMutex.Unlock()
	responseMacros[name] = macro
}

// CallMacro 调用已注册的响应宏
// 参数:
//   - name: 宏名称
//   - args: 传递给宏的参数
func CallMacro(name string, args ...interface{}) contracts.HttpResponse {
	macroMutex.RLock()
	macro, exists := responseMacros[name]
	macroMutex.RUnlock()

	if !exists {
		// 如果宏不存在，返回一个错误响应
		return Response("Macro not found: "+name, 500)
	}

	return macro(args...)
}

// HasMacro 检查是否存在指定名称的宏
func HasMacro(name string) bool {
	macroMutex.RLock()
	defer macroMutex.RUnlock()
	_, exists := responseMacros[name]
	return exists
}

// GetMacroNames 获取所有已注册的宏名称
func GetMacroNames() []string {
	macroMutex.RLock()
	defer macroMutex.RUnlock()

	names := make([]string, 0, len(responseMacros))
	for name := range responseMacros {
		names = append(names, name)
	}
	return names
}

// Api 示例宏函数（根据文档中的例子）
// 这个函数会在初始化时通过ResponseMacro注册
func Api(data interface{}, code ...int) contracts.HttpResponse {
	status := 200
	if len(code) > 0 {
		status = code[0]
	}

	return Json(map[string]interface{}{
		"data":    data,
		"status":  status,
		"message": "success",
	}, status)
}

// 初始化默认宏
func init() {
	// 注册默认的API宏
	ResponseMacro("api", func(args ...interface{}) contracts.HttpResponse {
		if len(args) == 0 {
			return Api(nil)
		}
		if len(args) == 1 {
			return Api(args[0])
		}
		if code, ok := args[1].(int); ok {
			return Api(args[0], code)
		}
		return Api(args[0])
	})
}