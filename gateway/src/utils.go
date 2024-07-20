package gateway

import (
	"fmt"
	"gateway/src/log"
	"gateway/src/pb"
	"github.com/golang/protobuf/proto"
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

func GetBinary(protoMsg proto.Message, log *log.Logger, agentConfig *AgentConfig) (res []byte, err error) {
	protoName := proto.MessageName(protoMsg)
	if agentConfig.EnableLogSend && protoName != "pb.s2cHeart" && protoName != "pb.s2cStrike" {
		log.Debug("send ==== %s %+v", protoName, protoMsg)
	}
	ba := pb.CreateEmpyByteArray()
	ba.WriteUint8(uint8(len(protoName)))
	ba.WriteString(protoName)
	binarybody, err := proto.Marshal(protoMsg)
	if err != nil {
		log.Error("proto.Marshal() failed: %+v", err)
		return nil, err
	}
	ba.WriteBytes(binarybody)
	return ba.Bytes(), nil
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

func RunOnAgentMgr(c chan Closure, mgr *AgentMgr, cb func(mgr *AgentMgr)) {
	c <- func() {
		cb(mgr)
	}
}

func RunOnAgent(c chan Closure, agent *Agent, cb func(agent *Agent)) {
	c <- func() {
		cb(agent)
	}
}

func RunOnMatch(c chan Closure, mgr *NatsMatch, cb func(mgr *NatsMatch)) {
	c <- func() {
		cb(mgr)
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
