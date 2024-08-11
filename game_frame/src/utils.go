package game_frame

import (
	"fmt"
	"game_frame/src/pb"
	"net"
	"strconv"
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

func RunOnRoom(c chan Closure, room *Room, cb func(room *Room)) {
	c <- func() {
		cb(room)
	}
}

// map[data: gameId:test4 gatewayIp:192.168.10.2 pbName:pb.LoginHallRequest roomId:test4_hall_1054 sn:2 timestamp:1.722344216e+09 uid:1054]
func ConvertRequest(dataMap map[string]interface{}) (*pb.CommonHead, []byte) {
	timestamp, _ := strconv.ParseInt(dataMap["timestamp"].(string), 10, 64)
	sn, _ := strconv.ParseInt(dataMap["sn"].(string), 10, 64)

	head := &pb.CommonHead{
		GameId:    dataMap["gameId"].(string),
		RoomId:    dataMap["roomId"].(string),
		Uid:       dataMap["uid"].(string),
		Pid:       dataMap["pid"].(string),
		Sn:        sn,
		Timestamp: timestamp,
		PbName:    dataMap["pbName"].(string),
		GatewayIp: dataMap["gatewayIp"].(string),
	}

	data := dataMap["data"].(string)

	return head, []byte(data)
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
