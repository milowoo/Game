package stresstest

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"golang.org/x/time/rate"

	"yalitest/stat"
)

// const
const (
	defaultClientCount      = 1
	defaultConnectPerSecond = 0
	defaultExecutePerSecond = 0
)

// Client ...
type Client interface {
	Connect() error
	Execute() error
	Close() error
}

// Manager ...
type Manager struct {
	// 生成的客户端数量 (default: 200)
	ClientCount int64
	// 每秒连接客户端数 (default: 0)
	// Note:
	//   * 0表示不限制
	ConnectPerSecond int64
	// 每个客户端每秒发起请求数 (default: 0)
	// Note:
	//   * 0表示不限制
	ExecutePerSecond int64
	// 生成客户端实例的函数
	ClientSpawnFunc func(n int64) Client
	// 压测报表数据
	Reporter *stat.Reporter

	startTm       time.Time
	connectStopTm time.Time
	wg            sync.WaitGroup
	stop          bool

	connectRequestCounter *stat.RequestCounter
	executeRequestCounter *stat.RequestCounter
}

// Start ...
func (m *Manager) Start() error {
	if m.ClientSpawnFunc == nil {
		return fmt.Errorf("ClientSpawnFunc must not be nil")
	}

	err := m.init()
	if err != nil {
		return err
	}

	connectRate := rate.Inf
	if m.ConnectPerSecond > 0 {
		connectRate = rate.Limit(m.ConnectPerSecond)
	}
	connectLimit := rate.NewLimiter(connectRate, int(connectRate))
	m.startTm = time.Now()

	for i := int64(0); i < m.ClientCount && !m.stop; i++ {
		_ = connectLimit.Wait(context.Background())

		c := m.ClientSpawnFunc(i)
		session := m.connectRequestCounter.NewRequest()
		err := c.Connect()
		if err != nil {
			session.Fail()
			_, _ = fmt.Fprintf(os.Stderr, "client connect fail: n=%v, err=%v", i, err)
			continue
		}

		session.Success()

		m.wg.Add(1)
		go m.doStress(c)
	}

	m.connectStopTm = time.Now()

	return nil
}

// Stop ...
func (m *Manager) Stop() {
	m.stop = true
}

// GracefulStop ...
func (m *Manager) GracefulStop() {
	m.Stop()
	m.Wait()
}

// Wait ...
func (m *Manager) Wait() {
	m.wg.Wait()
	fmt.Println(m.Report())
}

// Run ...
func (m *Manager) Run() error {
	err := m.Start()
	if err != nil {
		return err
	}

	m.Wait()
	return nil
}

// Report 输出压测报告
func (m *Manager) Report() string {
	now := time.Now()

	result := "\n"
	result += "********************************** stress test manager report **********************************\n"
	result += "\n"

	result += "压测配置: \n"
	result += fmt.Sprintf("客户端并发数: \t%v\n", m.ClientCount)
	result += fmt.Sprintf("每秒连接数: \t%v\n", m.ConnectPerSecond)
	result += fmt.Sprintf("每秒执行数: \t%v\n", m.ExecutePerSecond)

	result += "\n"
	result += fmt.Sprintf("持续时间: \t%v\n", now.Sub(m.startTm))

	connDuration := m.connectStopTm.Sub(m.startTm).Seconds()
	result += fmt.Sprintf("连接速率: \t%.2f [#/sec]\n", float64(m.connectRequestCounter.Total())/connDuration)

	execDuration := now.Sub(m.connectStopTm).Seconds()
	result += fmt.Sprintf("执行速率: \t%.2f [#/sec]\n", float64(m.executeRequestCounter.Total())/execDuration)

	if m.Reporter != nil {
		result += m.Reporter.String()
	}

	return result
}

func (m *Manager) init() error {
	if m.ClientCount <= 0 {
		m.ClientCount = defaultClientCount
	}

	if m.ConnectPerSecond < 0 {
		m.ConnectPerSecond = defaultConnectPerSecond
	}

	if m.ExecutePerSecond < 0 {
		m.ExecutePerSecond = defaultExecutePerSecond
	}

	m.connectRequestCounter = stat.NewRequestCounter(stat.NewTaskCounter(), stat.NewTimeRange(5000, 10))
	m.executeRequestCounter = stat.NewRequestCounter(stat.NewTaskCounter(), stat.NewTimeRange(5000, 10))

	if m.Reporter != nil {
		m.Reporter.AddCollector("connect", m.connectRequestCounter)
		m.Reporter.AddCollector("execute", m.executeRequestCounter)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT)
	go func() {
		<-sigChan
		m.Stop()
	}()

	return nil
}

func (m *Manager) doStress(c Client) {
	defer func() {
		p := recover()
		if p != nil {
			_, _ = fmt.Fprintf(os.Stderr, "execute panic recovered and going to stop: %v", p)
		}

		_ = c.Close()
		m.wg.Done()
	}()

	executeRate := rate.Inf
	if m.ExecutePerSecond > 0 {
		executeRate = rate.Limit(m.ExecutePerSecond)
	}
	executeLimit := rate.NewLimiter(executeRate, int(executeRate))

	for !m.stop {
		_ = executeLimit.Wait(context.Background())

		session := m.executeRequestCounter.NewRequest()
		err := c.Execute()
		if err != nil {
			session.Fail()
		} else {
			session.Success()
		}
	}
}
