package stat

import (
	"fmt"
	"sync"
)

// TimeRange ...
type TimeRange struct {
	Step        int64                             // 区间步长
	Count       int64                             // 区间数量
	UnitConvert func(v float64) (float64, string) // 时间单位转换

	maxVal     int64   // 最大值
	minVal     int64   // 最小值
	average    float64 // 平均数
	totalCount Counter

	data  []Counter
	mutex sync.Mutex
}

// NewTimeRange ...
func NewTimeRange(step, count int64) *TimeRange {
	result := &TimeRange{
		Step:        step,
		Count:       count,
		UnitConvert: MicroSecond,
	}
	for i := int64(0); i < count+1; i++ {
		result.data = append(result.data, Counter{})
	}
	return result
}

// Add ...
func (tm *TimeRange) Add(v int64) {
	if v > tm.maxVal {
		tm.maxVal = v
	}
	index := v / tm.Step

	if index >= tm.Count {
		index = tm.Count
	}

	tm.data[index].Add()

	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	curTotal := float64(tm.totalCount.Add())
	// 使用该算法避免时间总和溢出
	tm.average = tm.average*((curTotal-1)/curTotal) + float64(v)/curTotal
}

// Max ...
func (tm *TimeRange) Max() int64 {
	return tm.maxVal
}

// Min ...
func (tm *TimeRange) Min() int64 {
	return tm.minVal
}

// Average ...
func (tm *TimeRange) Average() float64 {
	return tm.average
}

// Summary ...
func (tm *TimeRange) Summary() string {
	sum := int64(0)
	for _, c := range tm.data {
		sum += c.Count
	}

	avg, avgUnit := tm.UnitConvert(tm.Average())
	max, maxUnit := tm.UnitConvert(float64(tm.Max()))
	min, minUnit := tm.UnitConvert(float64(tm.Min()))
	return fmt.Sprintf("<timeRange>: 总次数: %v, 平均延时: %.2f%v, 最大延时: %v%v, 最小延时: %v%v\n",
		sum, avg, avgUnit, max, maxUnit, min, minUnit)
}

// Detail ...
func (tm *TimeRange) Detail() string {
	result := "<timeRange>: \n"

	count := Counter{}

	for i := 0; i < len(tm.data)-1; i++ {
		count.AddCount(tm.data[i].Count)

		step, stepUnit := tm.UnitConvert(float64(tm.Step * int64(i+1)))
		result += fmt.Sprintf("    less than %.0f%v: \t%d\t(%.0f%% - %.0f%%)\n",
			step, stepUnit, tm.data[i].Count, count.PercentOf(&tm.totalCount), tm.data[i].PercentOf(&tm.totalCount))
	}

	last := int64(len(tm.data) - 1)
	lastStep, lastStepUnit := tm.UnitConvert(float64(tm.Step * last))
	result += fmt.Sprintf("    more than %.0f%v: \t%d\t(100%% - %.0f%%)\n",
		lastStep, lastStepUnit, tm.data[last].Count, tm.data[last].PercentOf(&tm.totalCount))

	avg, avgUnit := tm.UnitConvert(tm.Average())
	max, maxUnit := tm.UnitConvert(float64(tm.Max()))
	min, minUnit := tm.UnitConvert(float64(tm.Min()))

	result += fmt.Sprintf("    平均延迟:\t\t%.2f%v\n", avg, avgUnit)
	result += fmt.Sprintf("    最大延迟:\t\t%v%v\n", max, maxUnit)
	result += fmt.Sprintf("    最小延迟:\t\t%v%v\n", min, minUnit)
	result += fmt.Sprintf("    总次数:\t\t%d\n", tm.totalCount.Count)
	return result
}

// NanoSecond ...
func NanoSecond(v float64) (float64, string) {
	return v, "ns"
}

// NanoToMicroSecond 纳秒转换到微秒
func NanoToMicroSecond(v float64) (float64, string) {
	return v / 1000, "us"
}

// NanoToMilliSecond 纳秒转换到毫秒
func NanoToMilliSecond(v float64) (float64, string) {
	return v / 1000000, "ms"
}

// NanoToSecond 纳秒转换到秒
func NanoToSecond(v float64) (float64, string) {
	return v / 1000000000, "s"
}

// MicroSecond ...
func MicroSecond(v float64) (float64, string) {
	return v, "us"
}

// MicroToMilliSecond 微秒转换到毫秒
func MicroToMilliSecond(v float64) (float64, string) {
	return v / 1000, "ms"
}

// MicroToSecond 微秒转换到秒
func MicroToSecond(v float64) (float64, string) {
	return v / 1000000, "s"
}

// MilliSecond ...
func MilliSecond(v float64) (float64, string) {
	return v, "ms"
}

// MilliToSecond 毫秒转换到秒
func MilliToSecond(v float64) (float64, string) {
	return v / 1000, "s"
}
