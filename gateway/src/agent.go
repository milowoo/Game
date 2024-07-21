package gateway

import (
	"bytes"
	"compress/zlib"
	"gateway/src/constants"
	"gateway/src/handler"
	"gateway/src/log"
	"gateway/src/pb"
	"math"
	"net/url"
	"reflect"
	"strconv"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
)

const (
	// Maximum message size allowed from peer.
	MaxMsgSize              = 64 * 1024
	ClientPingInterval      = time.Millisecond * 2500
	ClientPingTimeout       = ClientPingInterval * 2
	ClientCheckLoginTimeout = ClientPingInterval * 4

	AgentFPS           = 20
	AgentFrameInterval = time.Second / AgentFPS

	CheckGameExceptionInternal = time.Second * 10 / AgentFPS

	SizeBytes = 0 // 对于websocket，没有sizebyte也ok
	MultiFlag = "pb.multi"
)

type SendEvent struct {
	t   int
	msg []byte
}

func NewSendBinaryEvent(msg []byte) *SendEvent {
	return &SendEvent{
		t:   websocket.BinaryMessage,
		msg: msg,
	}
}

type Agent struct {
	Server            *Server
	Config            *GlobalConfig
	DynamicConfig     *DynamicConfig
	Conn              *websocket.Conn
	SendQueue         chan *SendEvent
	RecvQueue         chan Closure
	MsgFromAgentMgr   chan Closure
	Closed            bool
	LastError         error
	AgentMgr          *AgentMgr
	LastPingTimeFrame int
	ConnectTimeFrame  int
	InHall            int //0 init 1 hall 2 room

	FrameTicker          *time.Ticker
	FrameID              int
	IsMatching           bool
	NextCheckPingFrame   int
	RequestGameErrFrame  int
	IsDisconnected       bool
	delayDisconnectTimer *time.Timer
	CachedMsgs           [][]byte
	Log                  *log.Logger
	Pid                  string
	Uid                  string
	RoomId               string
	GameId               string
	GameSubject          string
	Counter              *AtomicCounter
}

// new agent
func NewAgent(rawConn *websocket.Conn, agentMgr *AgentMgr, values url.Values) *Agent {
	timer := time.NewTimer(time.Hour)
	timer.Stop()

	log := agentMgr.Log
	self := &Agent{
		AgentMgr:             agentMgr,
		Server:               agentMgr.Server,
		DynamicConfig:        agentMgr.Server.DynamicConfig,
		Config:               agentMgr.Server.Config,
		Conn:                 rawConn,
		SendQueue:            make(chan *SendEvent, 2*1024),
		RecvQueue:            make(chan Closure, 2*1024),
		MsgFromAgentMgr:      make(chan Closure, 512),
		FrameID:              1,
		ConnectTimeFrame:     1,
		IsMatching:           true,
		InHall:               0,
		NextCheckPingFrame:   int(ClientPingInterval / AgentFrameInterval),
		delayDisconnectTimer: timer,
		CachedMsgs:           make([][]byte, 0, 64),
		Log:                  log,
		Uid:                  "",
		RoomId:               "",
		GameId:               "",
		GameSubject:          "",
		RequestGameErrFrame:  0,
		Counter:              &AtomicCounter{},
	}

	self.MsgFromAgentMgr <- func() {
		self.LoginGame(values.Get("game_id"), values.Get("timestamp"), values.Get("pid"), values.Get("token"))
	}
	return self
}

func (self *Agent) readPump() {
	defer func() {
		self.RecvQueue <- func() { self.onMultiClose() }
		close(self.RecvQueue)
		self.Conn.Close()

	}()

	for {
		t, msg, err := self.Conn.ReadMessage()
		if err != nil {
			if !websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseNoStatusReceived) {
				self.RecvQueue <- func() { self.onMultiError(err) }
			}

			break
		}

		self.LastPingTimeFrame = self.FrameID

		if t == websocket.TextMessage {
			self.RecvQueue <- func() { self.OnText(string(msg[:])) }
			continue
		}

		// t == websocket.BinaryMessage
		self.RecvQueue <- func() { self.OnBinary(msg) }
	}
}

func (self *Agent) writePump() {

	defer func() {
		self.Conn.Close()
	}()

	for {
		s, ok := <-self.SendQueue
		if !ok {
			// close channel.
			self.Conn.WriteMessage(websocket.CloseMessage, []byte{})
			break
		}

		self.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
		if err := self.Conn.WriteMessage(s.t, s.msg); err != nil {
			if !websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseNoStatusReceived) {
				self.Log.Error("uid %v WriteMessage() error, %+v", self.Uid, err)
			}

			break
		}
	}
}

func (self *Agent) onMultiClose() {
	if self.Closed {
		return
	}

	self.Conn.Close()
	self.Closed = true

	self.OnClose()
}

func (self *Agent) onMultiError(err error) {
	if self.LastError != nil {
		return
	}

	self.LastError = err
	self.OnError(err)
}

func (self *Agent) Run() {
	self.Conn.SetReadLimit(MaxMsgSize)
	self.Conn.SetPingHandler(func(string) error {
		self.RecvQueue <- func() { self.SendQueue <- &SendEvent{t: websocket.PongMessage, msg: []byte{}} }
		return nil
	})

	go self.readPump()
	go self.writePump()

	self.FrameTicker = time.NewTicker(AgentFrameInterval)
	defer func() {
		self.FrameTicker.Stop()
	}()

	self.OnOpen()

	for {
		select {
		case c, ok := <-self.RecvQueue:
			if !ok {
				return
			}
			SafeRunClosure(self, c)
		case c := <-self.MsgFromAgentMgr:
			SafeRunClosure(self, c)
		case <-self.FrameTicker.C:
			self.Frame()
		case <-self.delayDisconnectTimer.C:
			self.Log.Error("uid %v doDisconnect", self.Uid)
			self.onMultiClose()
			return
		}
	}
}

func (obj *Agent) SendBinary(msg []byte) {
	if !obj.Config.AgentConfig.EnableCachedMsg {
		obj.SendBinaryNow(msg)
		return
	}

	obj.CachedMsgs = append(obj.CachedMsgs, msg)
	if len(obj.CachedMsgs) >= obj.Config.AgentConfig.CachedMsgMaxCount {
		obj.FlushCachedMsgs()
		return
	}
}

func (self *Agent) SendBinaryNow(msg []byte) {
	if self.IsDisconnected || self.Closed {
		return
	}

	self.SendQueue <- NewSendBinaryEvent(msg)
}

func (self *Agent) OnOpen() {
	self.Log.Info("...")
}

/*
读取客户端的请求， 进行业务处理
*/
func (self *Agent) OnBinary(msg []byte) {
	ba := pb.CreateByteArray(msg)

	// 协议名长度
	dataLen, err := ba.ReadUint8()
	if err != nil {
		self.Log.Error("uid %v read proto head err %+v", self.Uid, err)
		return
	}

	// 协议名
	protoName, err := ba.ReadString(int(dataLen))
	if err != nil {
		self.Log.Error("uid %v read proto name err %+v", self.Uid, err)
		return
	}
	protoType := proto.MessageType(protoName)
	if protoType == nil {
		self.Log.Error("uid %v did not find proto ===== %s", self.Uid, protoName)
		return
	}

	// 协议结构体
	protoBody := reflect.New(protoType.Elem()).Interface().(proto.Message)
	pbBytes := make([]byte, ba.Available())
	_, err = ba.Read(pbBytes)
	if err != nil {
		self.Log.Error("uid %v read proto body err", self.Uid, err)
		return
	}
	err = proto.Unmarshal(pbBytes, protoBody)
	if err != nil {
		self.Log.Error("uid %v proto unmarshal err %+v", self.Uid, err)
		return
	}

	if self.Config.AgentConfig.EnableLogRecv && protoName != "pb.c2sHeart" && protoName != "pb.c2sStrike" {
		self.Log.Error("uid %+v receive ==== %s %+v", self.Uid, protoName, protoBody)
	}

	//检验游戏是否在开放
	if !self.checkGameOpen() {
		self.Log.Error("uid %v invalid game request", self.Uid)
		return
	}

	//如果不是心跳包，则需要转发给游戏服务
	if protoName == "pb.HeartbeatRequest" {
		handler.HeartBeatReply(self)
		return
	}

	if protoName == "pb.ClientLoginHallRequest" {
		handler.LoginHallRequest(self)
		return
	}

	if protoName == "pb.ClientMatchRequest" {
		handler.LoginHallRequest(self)
		return
	}
	if protoName == "pb.ClientMatchRequest" {
		//发起匹配请求
		handler.MatchRequest(self)
		return
	}

	if protoName == "pb.ClientCancelMatchRequest" {
		//发起取消匹配请求
		handler.CancelMatchRequest(self)
		return
	}

	//转发消息给游戏服务器，并接收应答给客户端
	handler.ForwardClientRequest(self, protoName, protoBody)

}

func (self *Agent) checkGameOpen() bool {
	gameInfo := self.DynamicConfig.GetGameInfo(self.GameId)
	if gameInfo == nil {
		self.Log.Warn("LoginHallRequest uid %+v get game info err gameId %+v", self.Uid, self.GameId)
		return false
	}

	if gameInfo.Status != 1 {
		self.Log.Warn("LoginHallRequest uid %+v game status %+v err gameId %+v", self.Uid, gameInfo.Status, self.GameId)
		self.DelayDisconnect(5)
		return false
	}

	return true
}

func (self *Agent) OnText(msg string) {
	self.Log.Info("on longer supports text message %s", msg)
}

func (self *Agent) OnError(err error) {
	self.Log.Error("uid %v %+v", self.Uid, err)
}

func (self *Agent) OnClose() {
	self.Log.Info("OnClose uid %v ...", self.Uid)

	close(self.SendQueue)

	if !self.IsDisconnected {
		self.notifyLostAgent(true)
	}
}

func (self *Agent) CompresssData(data []byte) []byte {
	var b bytes.Buffer
	w := zlib.NewWriter(&b)
	w.Write(data)
	w.Close()

	compressedData := b.Bytes()
	compressedDataLen := len(compressedData)
	sendData := make([]byte, 0, 1+SizeBytes+compressedDataLen)
	sendData = append(sendData, 'c')

	if SizeBytes >= 2 {
		sendData = append(sendData, byte(compressedDataLen&0xff)) // 小端
		sendData = append(sendData, byte((compressedDataLen>>8)&0xff))
	}

	if SizeBytes == 4 {
		sendData = append(sendData, byte((compressedDataLen>>16)&0xff))
		sendData = append(sendData, byte((compressedDataLen>>24)&0xff))
	}

	sendData = append(sendData, compressedData...)

	return sendData
}

func (self *Agent) FlushCachedMsgs() {
	msgCount := len(self.CachedMsgs)
	if msgCount <= 0 {
		return
	}

	if msgCount == 1 {
		data := self.CachedMsgs[0]
		self.SendBinaryNow(data)
	} else {
		ba := pb.CreateEmpyByteArray()
		ba.WriteUint8(uint8(len(MultiFlag)))
		ba.WriteString(MultiFlag)
		for _, binary := range self.CachedMsgs {
			ba.WriteInt32(int32(len(binary)))
			ba.WriteBytes(binary)
		}
		self.SendBinaryNow(ba.Bytes())
	}
	self.CachedMsgs = self.CachedMsgs[:0]
	return

}

func (self *Agent) notifyLostAgent(noticeGame bool) {
	if len(self.Uid) <= 1 {
		return
	}

	//通知agentmgr 删除 agent信息
	RunOnAgentMgr(self.AgentMgr.MsgFromAgent, self.AgentMgr, func(agentMgr *AgentMgr) {
		agentMgr.loseConnect(self)
	})

	//关闭链接
	self.CloseConnect()

	//通知 game_x 用户失去链接
	if noticeGame {
		handler.UserExitHandler(self, 1)
	}

}

func (self *Agent) CloseConnect() {
	if self.Closed {
		return
	}

	self.Conn.Close()
	self.Closed = true
}

func (self *Agent) checkInvalidAgent() {
	pingElapseTime := time.Duration(self.FrameID-self.ConnectTimeFrame) * AgentFrameInterval
	if pingElapseTime < ClientCheckLoginTimeout {
		return
	}

	self.notifyLostAgent(true)

	self.Log.Info("checkInvalidAgent uid %v delayDisconnect: ping timeout", self.Uid)
	self.DelayDisconnect(0)
}

func (self *Agent) CheckPing() {
	if !self.Config.AgentConfig.EnableCheckPing {
		return
	}

	pingElapseTime := time.Duration(self.FrameID-self.LastPingTimeFrame) * AgentFrameInterval
	if pingElapseTime < ClientPingTimeout {
		return
	}

	self.notifyLostAgent(true)
	self.Log.Info("CheckPing uid %v delayDisconnect: ping timeout", self.Uid)
	self.DelayDisconnect(0)
}

func (self *Agent) Frame() {
	self.FrameID++

	if self.FrameID >= self.NextCheckPingFrame {
		self.NextCheckPingFrame = self.FrameID + int(ClientPingInterval/AgentFrameInterval)
		self.CheckPing()
	}

	//校验是否是非法的链接
	self.checkInvalidAgent()

	self.FlushCachedMsgs()

	self.checkGameException()
}

func (self *Agent) checkGameException() {
	if self.RequestGameErrFrame == 0 {
		return
	}
	if self.RequestGameErrFrame < self.FrameID-int(CheckGameExceptionInternal/AgentFrameInterval) {
		//如果是在大厅，则需要游戏服创建大厅
		if self.InHall == 1 {
			handler.LonginHall2Match(self)
			return
		}
		self.notifyLostAgent(false)

		self.Log.Info("checkInvalidAgent uid %v delayDisconnect: ping timeout", self.Uid)
		self.DelayDisconnect(0)
	}
}

func (self *Agent) DelayDisconnect(delay time.Duration) {
	self.Log.Info("uid %v try delayDisconnect, %.02fs", self.Uid, delay.Seconds())

	self.IsDisconnected = true
	self.delayDisconnectTimer.Reset(delay)
}

func (self *Agent) LoginGame(gameId, t, pid, token string) int {
	if self.Config.AgentConfig.EnableLogRecv {
		self.Log.Debug(" receive ====  %+v %+s %s %s %v", t, gameId, pid, token)
	}

	timestamp, err := strconv.ParseInt(t, 10, 64)
	if err != nil {
		handler.DoLoginReply(self, constants.INVALID_TIME, "unexpected timestamp", 0)
		return -1
	}

	config := self.Config.AgentConfig
	if config.EnableCheckLoginParams && math.Abs(float64(time.Now().Sub(time.Unix(timestamp, 0)))) >= float64(config.TimestampExpireDuration) {
		handler.DoLoginReply(self, constants.EXPIRED_TIME, "timestamp expired", 0)
		return -1
	}

	gameInfo := self.DynamicConfig.GetGameInfo(gameId)
	if gameInfo == nil {
		handler.DoLoginReply(self, constants.INVALID_GAME_ID, "invalid gameId", 0)
		return -1
	}

	if gameInfo.Status != 1 {
		handler.DoLoginReply(self, constants.GAME_IS_OFF, "game offline", 0)
		return -1
	}

	//去平台验签

	//调用ucenter获取 userId
	uid, err := g_Server.UCenterMgr.ApplyUid(pid)
	if err != nil {
		self.Log.Error("EnterGame  %+v apply err %+v", pid, err)
		handler.DoLoginReply(self, constants.SYSTEM_ERROR, "system err", 0)
		return -1
	}
	self.Pid = pid
	self.Uid = uid

	//反射到agentMgr
	RunOnAgentMgr(self.AgentMgr.MsgFromAgent, self.AgentMgr, func(agentMgr *AgentMgr) {
		agentMgr.EnterGame(self)
	})

	handler.DoLoginReply(self, 0, "success", 0)
	return 0
}

func (self *Agent) MatchResponse(res *pb.MatchOverRes) {
	handler.MatchOverResponse(self, res)
}

func (self *Agent) ProcGamePushMessage(res *pb.GamePushMessage) {
	head := res.GetHead()
	self.Log.Info("ProcGamePushMessage gameId %+v uid %+v protoName %+v ", head.GetGameId(), head.GetUid(), head.GetProtoName())
	var protoMessage proto.Message
	proto.Unmarshal(res.GetData(), protoMessage)
	self.ReplyClient(protoMessage)
}

func (self *Agent) ReplyClient(protoMsg proto.Message) {
	binary, err := GetBinary(protoMsg, self.Log, self.Config.AgentConfig)
	if err != nil {
		return
	}
	self.SendBinaryNow(binary)
}
