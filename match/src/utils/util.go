package utils

import (
	"encoding/json"
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/v2/model"
	"math/rand"
	"reflect"
	"sync/atomic"
)

func JsonToServices(result string) (*[]model.Service, error) {
	var services []model.Service
	err := json.Unmarshal([]byte(result), &services)
	if err != nil {
		return nil, err
	}
	//if len(service.Hosts) == 0 {
	//	logger.Warnf("instance list is empty,json string:%s", result)
	//}
	return &services, nil

}
func ToJsonString(object interface{}) string {
	js, _ := json.Marshal(object)
	return string(js)
}

func ConvertInterfaceToString(data interface{}) (string, error) {
	// 使用 reflect 包检查 data 是否为 string 类型
	if reflect.TypeOf(data).Kind() != reflect.String {
		return "", fmt.Errorf("expected a string, got %T", data)
	}

	// 如果是 string 类型，返回其数据
	return data.(string), nil
}

type Closure = func()

func SafeRunClosure(v interface{}, c Closure) {
	defer func() {
		if err := recover(); err != nil {
			//log.Printf("%+v: %s", err, debug.Stack())

		}
	}()

	c()
}

func RandomInt(rand *rand.Rand, min int64, maxPlus1 int64) int64 {
	return min + rand.Int63n((maxPlus1 - min))
}

// AtomicCounter 是一个使用原子操作实现的线程安全计数器
type AtomicCounter struct {
	value int64
}

// Increment 原子递增计数器
func (c *AtomicCounter) Increment() {
	atomic.AddInt64(&c.value, 1)
}

// Increment 原子递增计数器
func (c *AtomicCounter) GetIncrementValue() int64 {
	atomic.AddInt64(&c.value, 1)
	return atomic.LoadInt64(&c.value)
}

// Value 获取计数器的当前值
func (c *AtomicCounter) Value() int64 {
	return atomic.LoadInt64(&c.value)
}
