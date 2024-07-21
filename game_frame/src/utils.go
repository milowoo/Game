package game_frame

import (
	"fmt"
	"game_frame/src/handler"
	"net"
	"sync/atomic"
)

func GetHostIp() string {
	addrList, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println("get current host ip err: ", err)
		return ""
	}
	var ip string
	for _, address := range addrList {
		if ipNet, ok := address.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				ip = ipNet.IP.String()
				break
			}
		}
	}
	return ip
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

func RunOnRoomMgr(c chan Closure, mgr *RoomMgr, cb func(mgr *RoomMgr)) {
	c <- func() {
		cb(mgr)
	}
}

func RunOnRoom(c chan Closure, room *handler.Room, cb func(room *handler.Room)) {
	c <- func() {
		cb(room)
	}
}

func RunOnNatsService(c chan Closure, natsService *NatsService, cb func(natsService *NatsService)) {
	c <- func() {
		cb(natsService)
	}
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
