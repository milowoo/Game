package mq

import (
	"fmt"
	"github.com/nats-io/go-nats"
	"sync"
	"time"
)

// nats连接池值类型
type NatsPool struct {
	conns    chan *nats.Conn
	dialFunc DialFunc
	stopOnce sync.Once
	Network  string
	Addr     string
}

// 连接处理函数
type DialFunc func(natsServersUrl string, options ...nats.Option) (*nats.Conn, error)

// 创建连接池的工厂方法
func NewNatsConnectPool(addr string, connSize int, dialFunc DialFunc) (*NatsPool, error) {
	var conn *nats.Conn
	var err error
	pool := make([]*nats.Conn, 0, connSize)
	for i := 0; i < connSize; i++ {
		conn, err = dialFunc(addr)
		if err != nil {
			for _, conn = range pool {
				conn.Close()
			}
			pool = pool[:0]
			break
		}
		pool = append(pool, conn)
	}
	p := NatsPool{
		Addr:     addr,
		conns:    make(chan *nats.Conn, len(pool)),
		dialFunc: dialFunc,
	}
	for i := range pool {
		p.conns <- pool[i]
	}

	if connSize < 1 {
		return &p, err
	}

	return &p, err
}

const (
	DefaultConnSize = 20 // 默认初始的连接数
)

// 默认的连接处理函数
var DefaultDialFunc = func(natsServersUrl string, options ...nats.Option) (*nats.Conn, error) {
	ops := []nats.Option{
		nats.MaxReconnects(5),               // 设置重新连接等待和最大重新连接尝试次数
		nats.ReconnectWait(2 * time.Second), // 每次重连等待时间
		//nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
		//	// fmt.Println("Nats server disconnected Reason:" + err.Error())
		//	logger.WarmLog("Nats server disconnected Reason:" + err.Error())
		//}), // 断开连接的错误处理
		nats.ReconnectHandler(func(nc *nats.Conn) {
			fmt.Println("Nats server reconnected to " + nc.ConnectedUrl())
		}), // 重连时的错误处理
		nats.ClosedHandler(func(nc *nats.Conn) {
			fmt.Println("Nats server connection closed. Reason: " + nc.LastError().Error())
		}), // 关闭连接时的错误处理
	}

	ops = append(ops, options...)

	return nats.Connect(natsServersUrl, ops...)
}

// 默认连接池的工厂方法
func NewDefaultPool(addr string) (*NatsPool, error) {
	return NewNatsConnectPool(addr, DefaultConnSize, DefaultDialFunc)
}

// 从连接池获取连接
func (p *NatsPool) Get() (*nats.Conn, error) {
	select {
	case conn := <-p.conns:
		return conn, nil
	default:
		return p.dialFunc(p.Addr)
	}
}

// 将连接放回连接池
func (p *NatsPool) Put(conn *nats.Conn) {
	select {
	case p.conns <- conn:
	default:
		conn.Close()
	}
}

// 情况连接池
func (p *NatsPool) Empty() {
	var conn *nats.Conn
	for {
		select {
		case conn = <-p.conns:
			conn.Close()
		default:
			return
		}
	}
}

// 有效的连接数
func (p *NatsPool) Avail() int {
	return len(p.conns)
}
