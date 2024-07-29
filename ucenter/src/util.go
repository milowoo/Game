package ucenter

import (
	"fmt"
	"reflect"
)

type Closure = func()

func ConvertInterfaceToString(data interface{}) (string, error) {
	// 使用 reflect 包检查 data 是否为 string 类型
	if reflect.TypeOf(data).Kind() != reflect.String {
		return "", fmt.Errorf("expected a string, got %T", data)
	}

	// 如果是 string 类型，返回其数据
	return data.(string), nil
}
