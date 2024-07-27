package stat

import (
	"fmt"
	"sync/atomic"
)

// Counter ...
type Counter struct {
	Count int64
}

// NewCounter ...
func NewCounter() *Counter {
	return &Counter{}
}

// Add ...
func (c *Counter) Add() int64 {
	return atomic.AddInt64(&c.Count, 1)
}

// AddCount ...
func (c *Counter) AddCount(count int64) int64 {
	return atomic.AddInt64(&c.Count, count)
}

// RateOf ...
func (c *Counter) RateOf(total *Counter) float64 {
	if total.Count == 0 {
		return 0
	}
	return float64(c.Count) / float64(total.Count)
}

// PercentOf ...
func (c *Counter) PercentOf(total *Counter) float64 {
	if total.Count == 0 {
		return 0
	}
	return float64(c.Count) / float64(total.Count) * 100
}

// Summary ...
func (c *Counter) Summary() string {
	return fmt.Sprintf("<counter>: 次数: %d\n", c.Count)
}
