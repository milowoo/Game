package util

import (
	"github.com/google/uuid"
	"math/rand"
	"sync/atomic"
	"time"
)

func RandomInt(rand *rand.Rand, min int, maxPlus1 int) int {
	return min + int(rand.Int63n(int64(maxPlus1-min)))
}

func UUID() string {
	return uuid.New().String()
}

// 获取时间戳(毫秒）
func getCurDay() string {
	nTime := time.Now()
	curDay := nTime.Format("20060102")
	return curDay
}

func getPreDay() string {
	nTime := time.Now()
	preTime := nTime.AddDate(0, 0, -1)
	preDay := preTime.Format("20060102")
	return preDay
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
