package main

import (
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"log"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"test"
	log2 "test/log"
	"test/pb"
	"time"
)

var gCounter *test.AtomicCounter
var gLog *log2.Logger
var g_gameId = "test4"
var g_isInHall = false
var g_Uid = "100"

type URL struct {
	// 连接地址
	Addr string
	// 游戏名称 (game id)
	GameName string
	// 厂商AppKey
	AppKey string
	// uid
	UID string
}

func main() {
	logName := "/Users/wuchuangeng/game/logs/" + "new_test.log"
	Log := log2.NewLogger2(logName, 1024*2, log.LstdFlags|log.Lshortfile, 8, 1024, 2)
	gLog = Log

	gCounter = &test.AtomicCounter{}
	gLog.Info("test begin .........")

	u := &URL{
		Addr:     "ws://127.0.0.1:36001",
		GameName: g_gameId,
	}

	// 解析URL并创建WebSocket连接
	log.Printf("Connecting to WebSocket server at %s", u.String())

	// 创建一个新的WebSocket连接
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("Dial failed: ", err)
	}
	defer conn.Close()

	// 注册中断信号处理器
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	// 接收消息的goroutine
	go func() {
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				gLog.Error("Read error:", err)
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					gLog.Error("Unexpected close error, exiting...")
					return
				}
				return
			}
			Receive(message)
		}
	}()

	// 发送消息的goroutine
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				Execute(conn)
				HeartBeat(conn)
			}
		}
	}()

	// 等待中断信号
	<-interrupt

	log.Println("Exiting...")
}

func Receive(message []byte) {
	var response pb.ClientCommonResponse
	proto.Unmarshal(message, &response)
	gLog.Info("Receive %+v", response)
}

func (u *URL) String() string {
	tm := time.Now().Unix()

	rawURL, _ := url.Parse(u.Addr)

	v := url.Values{}
	v.Add("timestamp", strconv.FormatInt(tm, 10))
	v.Add("token", test.UUID())
	v.Add("pid", test.UUID())
	v.Add("gameId", u.GameName)
	rawURL.RawQuery = v.Encode()

	return rawURL.String()
}

func HeartBeat(conn *websocket.Conn) {
	req := &pb.HeartbeatRequest{
		Timestamp: time.Now().Unix(),
	}
	bytes, _ := proto.Marshal(req)
	head := &pb.ClientCommonHead{
		Pid:       test.UUID(),
		Sn:        gCounter.GetIncrementValue(),
		ProtoName: proto.MessageName(req),
		Timestamp: time.Now().Unix(),
	}

	request := &pb.ClientCommonRequest{
		Head: head,
		Body: bytes,
	}

	d, err := proto.Marshal(request)

	err = conn.WriteMessage(websocket.BinaryMessage, d)
	if err != nil {
		gLog.Error("Write error:", err)
		return
	}

	gLog.Info("message client send  heartbeat success")
}

func Execute(conn *websocket.Conn) {
	SendEndHall(conn)
}

func SendEndHall(conn *websocket.Conn) {
	if g_isInHall {
		return
	}

	g_isInHall = true

	gLog.Info("send end hall uid %+v", g_Uid)
	enterRoomReq := &pb.ClientLoginHallRequest{}
	bytes, _ := proto.Marshal(enterRoomReq)

	head := &pb.ClientCommonHead{
		Pid:       g_Uid,
		Sn:        gCounter.GetIncrementValue(),
		ProtoName: proto.MessageName(enterRoomReq),
		Timestamp: time.Now().Unix(),
	}

	request := &pb.ClientCommonRequest{
		Head: head,
		Body: bytes,
	}

	d, err := proto.Marshal(request)

	err = conn.WriteMessage(websocket.BinaryMessage, d)
	if err != nil {
		gLog.Error("Write error:", err)
		return
	}

	gLog.Info("SendEndHall end hall success")

	g_isInHall = true
	return
}
