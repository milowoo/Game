package stat

import (
	"fmt"
)

// TaskCounter ...
type TaskCounter struct {
	Total   Counter
	Success Counter
	Fail    Counter
}

// NewTaskCounter ...
func NewTaskCounter() *TaskCounter {
	return &TaskCounter{}
}

// AddTotal ...
func (tc *TaskCounter) AddTotal() int64 {
	return tc.Total.Add()
}

// AddTotalCount ...
func (tc *TaskCounter) AddTotalCount(n int64) int64 {
	return tc.Total.AddCount(n)
}

// AddSuccess ...
func (tc *TaskCounter) AddSuccess() int64 {
	return tc.Success.Add()
}

// AddSuccessCount ...
func (tc *TaskCounter) AddSuccessCount(n int64) int64 {
	return tc.Success.AddCount(n)
}

// AddFail ...
func (tc *TaskCounter) AddFail() int64 {
	return tc.Fail.Add()
}

// AddFailCount ...
func (tc *TaskCounter) AddFailCount(n int64) int64 {
	return tc.Fail.AddCount(n)
}

// Summary ...
func (tc *TaskCounter) Summary() string {
	result := ""

	missCount := tc.Total.Count - tc.Success.Count - tc.Fail.Count
	missPercent := float64(missCount) / float64(tc.Total.Count) * 100

	result += fmt.Sprintf("<task>: 总次数: %v, 成功数: %v(%.2f%%), 失败数: %v(%.2f%%), 丢失数: %v(%.2f%%)\n",
		tc.Total.Count, tc.Success.Count, tc.Success.PercentOf(&tc.Total),
		tc.Fail.Count, tc.Fail.PercentOf(&tc.Total), missCount, missPercent)

	return result
}
