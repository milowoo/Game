package stat

import (
	"fmt"
	"time"
)

// RequestCounter ...
type RequestCounter struct {
	taskCounter *TaskCounter
	timeRange   *TimeRange
}

// NewRequestCounter ...
func NewRequestCounter(taskCounter *TaskCounter, timeRange *TimeRange) *RequestCounter {
	return &RequestCounter{
		taskCounter: taskCounter,
		timeRange:   timeRange,
	}
}

// Total ...
func (rc *RequestCounter) Total() int64 {
	return rc.taskCounter.Total.Count
}

// Success ...
func (rc *RequestCounter) Success() int64 {
	return rc.taskCounter.Success.Count
}

// Fail ...
func (rc *RequestCounter) Fail() int64 {
	return rc.taskCounter.Fail.Count
}

// NewRequest ...
func (rc *RequestCounter) NewRequest() *RequestSession {
	rc.taskCounter.AddTotal()
	return &RequestSession{
		parent:  rc,
		startTm: time.Now(),
	}
}

// Summary ...
func (rc *RequestCounter) Summary() string {
	result := "<request>: \n"
	result += fmt.Sprintf("  * 次数统计 - %v", rc.taskCounter.Summary())
	result += fmt.Sprintf("  * 延迟统计 - %v", rc.timeRange.Summary())
	return result
}

// Detail ...
func (rc *RequestCounter) Detail() string {
	result := "<request>: **** begin ****\n"
	result += rc.timeRange.Detail()
	result += "<request>: **** end ****\n"
	return result
}

// RequestSession ...
type RequestSession struct {
	parent  *RequestCounter
	startTm time.Time
}

// Success ...
func (rs *RequestSession) Success() {
	rs.parent.taskCounter.AddSuccess()
	rs.parent.timeRange.Add(time.Now().Sub(rs.startTm).Microseconds())
}

// Fail ...
func (rs *RequestSession) Fail() {
	rs.parent.taskCounter.AddFail()
	rs.parent.timeRange.Add(time.Now().Sub(rs.startTm).Microseconds())
}
