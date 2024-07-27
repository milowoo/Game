package stat

import (
	"fmt"
	"sync"
	"time"
)

// AsyncRequestCounter ...
type AsyncRequestCounter struct {
	taskCounter *TaskCounter
	timeRange   *TimeRange

	data  map[interface{}]time.Time
	mutex sync.Mutex
}

// NewAsyncRequestCounter ...
func NewAsyncRequestCounter(taskCounter *TaskCounter, timeRange *TimeRange) *AsyncRequestCounter {
	return &AsyncRequestCounter{
		taskCounter: taskCounter,
		timeRange:   timeRange,
		data:        map[interface{}]time.Time{},
	}
}

// NewAsyncRequestCounterDefault ...
func NewAsyncRequestCounterDefault() *AsyncRequestCounter {
	return &AsyncRequestCounter{
		taskCounter: NewTaskCounter(),
		timeRange:   NewTimeRange(5000, 10),
		data:        map[interface{}]time.Time{},
	}
}

// AddRequest ...
func (arc *AsyncRequestCounter) AddRequest(key interface{}) {
	if key == nil {
		return
	}

	arc.taskCounter.AddTotal()

	arc.mutex.Lock()
	defer arc.mutex.Unlock()

	if arc.data == nil {
		arc.data = map[interface{}]time.Time{}
	}

	arc.data[key] = time.Now()
}

// AddSuccess ...
func (arc *AsyncRequestCounter) AddSuccess(key interface{}) {
	arc.addResponse(key, arc.taskCounter.AddSuccess)
}

// AddFail ...
func (arc *AsyncRequestCounter) AddFail(key interface{}) {
	arc.addResponse(key, arc.taskCounter.AddFail)
}

func (arc *AsyncRequestCounter) addResponse(key interface{}, fn func() int64) {
	if key == nil {
		return
	}

	arc.mutex.Lock()
	defer arc.mutex.Unlock()

	if arc.data == nil {
		return
	}

	startTm, exists := arc.data[key]
	if !exists {
		return
	}

	fn()
	arc.timeRange.Add(time.Now().Sub(startTm).Microseconds())

	delete(arc.data, key)
}

// Summary ...
func (arc *AsyncRequestCounter) Summary() string {
	result := "<asyncRequest>: \n"
	result += fmt.Sprintf("  * 次数统计 - %v", arc.taskCounter.Summary())
	result += fmt.Sprintf("  * 延迟统计 - %v", arc.timeRange.Summary())
	return result
}

// Detail ...
func (arc *AsyncRequestCounter) Detail() string {
	result := "<asyncRequest>: **** begin ****\n"
	result += arc.timeRange.Detail()
	result += "<asyncRequest>: **** end ****\n"
	return result
}
