package stat

import (
	"fmt"
)

// Collector 数据收集器
type Collector interface {
	Summary() string
}

// DetailCollector 拥有详细报告的数据收集器
type DetailCollector interface {
	Collector
	Detail() string
}

type namedCollector struct {
	name      string
	collector Collector
}

// Reporter ...
type Reporter struct {
	collectors []namedCollector
}

// NewReporter ...
func NewReporter() *Reporter {
	return &Reporter{
		collectors: make([]namedCollector, 0, 32),
	}
}

// String ...
func (t *Reporter) String() string {
	result := "\n"
	result += "********************************** Reporter **********************************\n"
	result += "\n"

	result += "* 简要统计信息: \n"
	for _, c := range t.collectors {
		result += fmt.Sprintf("%v -> %v\n", c.name, c.collector.Summary())
	}

	result += "\n"
	result += "* 详细统计信息: \n"
	for _, c := range t.collectors {
		dc, ok := c.collector.(DetailCollector)
		if !ok {
			continue
		}

		result += fmt.Sprintf("%v -> %v\n", c.name, dc.Detail())
	}

	result += "\n"
	return result
}

// AddCollector 增加采集器数据到压测报告中
func (t *Reporter) AddCollector(name string, collector Collector) {
	t.collectors = append(t.collectors, namedCollector{
		name:      name,
		collector: collector,
	})
}
