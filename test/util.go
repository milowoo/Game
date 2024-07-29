package test

import (
	"github.com/google/uuid"
	"sync/atomic"
)

func UUID() string {
	return uuid.New().String()
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
