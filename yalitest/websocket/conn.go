package websocket

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"os"
	"sync"
	"yalitest/log"
	"yalitest/pb"

	"github.com/gorilla/websocket"
)

// Connection ...
type Connection struct {
	Module            string
	OnDownstream      func(response *pb.ClientCommonResponse)
	OnConnectionError func(e error)

	Log *log.Logger

	seq   int64
	conn  *websocket.Conn
	mutex sync.Mutex
}

// Dial ...
func (c *Connection) Dial(addr string) error {
	c.Log.Info(" conn dial addr %+v", addr)
	dialer := websocket.Dialer{}
	conn, _, err := dialer.Dial(addr, nil)
	if err != nil {
		return fmt.Errorf("channel dial conn fail: %v", err)
	}

	c.conn = conn

	go c.serve()
	return nil
}

// Upstream ...
func (c *Connection) Upstream(message proto.Message) error {
	d, err := proto.Marshal(message)
	if err != nil {
		return fmt.Errorf("channel marshal upstream pb fail: %v", err)
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	err = c.conn.WriteMessage(websocket.BinaryMessage, d)
	if err != nil {
		return err
	}

	return nil
}

// Close ...
func (c *Connection) Close() {
	_ = c.conn.Close()
}

func (c *Connection) serve() {
	for {
		err := c.doServe()
		if err != nil {
			break
		}
	}
}

func (c *Connection) doServe() error {
	defer func() {
		p := recover()
		if p != nil {
			_, _ = fmt.Fprintf(os.Stderr, "handle downstream panic recovered: %v", p)
		}
	}()

	_, d, err := c.conn.ReadMessage()
	if err != nil {
		if c.OnConnectionError != nil {
			c.OnConnectionError(err)
		}
		return err
	}

	var response pb.ClientCommonResponse
	proto.Unmarshal(d, &response)

	if c.OnDownstream != nil {
		c.OnDownstream(&response)
	}

	return nil
}
