package game_frame

import (
	"game_frame/src/constants"
	"game_frame/src/log"
	"game_frame/src/pb"
	"game_frame/src/redis"
	"github.com/golang/protobuf/proto"
	"math/rand"
	"sync"
	"time"
)

const (
	RoomMgrFrameInterval = time.Millisecond * 500
)

type RoomMgr struct {
	Server      *Server
	Log         *log.Logger
	MsgFromNats chan Closure
	MsgFromRoom chan Closure

	GlobalConfig  *GlobalConfig
	DynamicConfig *DynamicConfig

	RedisDao *redis.RedisDao

	rand      *rand.Rand
	roomMutex sync.Mutex

	frameTimer            *time.Ticker
	id2Room               map[string]*Room
	innerId2Room          map[string]*Room
	innerId2RoomId        map[string]string
	roomId2InnerId        map[string]string
	frameId               int
	nextLogRoomCountFrame int

	isQuit bool
	exit   chan bool
}

func NewRoomMgr(sever *Server) *RoomMgr {
	now := time.Now()
	roomMgr := &RoomMgr{
		Server:        sever,
		MsgFromNats:   make(chan Closure, 10*1024),
		MsgFromRoom:   make(chan Closure, 2*1024),
		GlobalConfig:  sever.Config,
		DynamicConfig: sever.DynamicConfig,
		RedisDao:      sever.RedisDao,

		rand:       rand.New(rand.NewSource(now.Unix())),
		frameTimer: time.NewTicker(RoomMgrFrameInterval),
		id2Room:    make(map[string]*Room, 0),

		frameId:               0,
		nextLogRoomCountFrame: 0,
		isQuit:                false,
		exit:                  make(chan bool, 1),
	}

	return roomMgr
}

func (self *RoomMgr) Run() {
	defer func() {
		p := recover()
		if p != nil {
			self.Log.Info("execute panic recovered and going to stop: %v", p)
		}
	}()

	self.Server.WaitGroup.Add(1)
	defer func() {
		self.Server.WaitGroup.Done()
	}()

	for {
		// 优先查看exit，
		select {
		case <-self.exit:
			{
				self.Quit()
				return
			}
		case c := <-self.MsgFromNats:
			SafeRunClosure(self, c)
		case c := <-self.MsgFromRoom:
			SafeRunClosure(self, c)
		case <-self.frameTimer.C:
			SafeRunClosure(self, func() {
				self.Frame()
			})
		}
	}

}

func (self *RoomMgr) Frame() {

}

func (self *RoomMgr) Quit() {
	// 通知所有room强制存储并退出
	for _, room := range self.id2Room {
		RunOnRoom(room.MsgFromMgr, room, func(input *Room) {
			input.SaveAndQuit()
		})
	}

	for len(self.id2Room) > 0 {
		c := <-self.MsgFromRoom
		SafeRunClosure(self, c)
	}
}

func (self *RoomMgr) ProcessGatewayRequest(reply string, msg interface{}) {
	self.Log.Info("process gateway request begin .....")
	var commonReq pb.GameCommonRequest
	commonRes := &pb.GameCommonResponse{}

	req, _ := ConvertInterfaceToString(msg)
	proto.Unmarshal([]byte(req), &commonReq)
	head := commonReq.GetHead()
	if head.GetGameId() != self.Server.Config.GameId {
		self.Log.Error("ProcessGatewayRequest invalid nat data %+v ", commonReq.GetHead().GameId)
		commonRes.Code = constants.INVALID_BODY
		commonRes.Msg = "Unmarshal request err"
		res, _ := proto.Marshal(commonRes)
		self.Server.NatsPool.Publish(reply, map[string]interface{}{"res": "ok", "data": string(res)})
		return
	}

	var request proto.Message
	proto.Unmarshal(commonReq.GetData(), request)

	room := self.GetOrCreateRoom(head)
	if room == nil {
		self.Log.Error("ProcessGatewayRequest invalid proto %+v gameId %+v uid %+v",
			head.ProtoName, head.GetRoomId(), head.GetUid())
		commonRes.Code = constants.SYSTEM_ERROR
		commonRes.Msg = "system err"
		res, _ := proto.Marshal(commonRes)
		self.Server.NatsPool.Publish(reply, map[string]interface{}{"res": "ok", "data": string(res)})
		return
	}

	RunOnRoom(room.MsgFromMgr, room, func(input *Room) {
		room.ApplyProtoHandler(reply, head, request)
	})

	return
}

func (self *RoomMgr) GetOrCreateRoom(head *pb.CommonHead) *Room {
	//先判断房间是否在使用中
	room, ok := self.id2Room[head.RoomId]
	if ok {
		return room
	}

	//如果不是进入房间/ 大厅的协议，就是非法的请求
	if head.ProtoName != "pb.LoginHallRequest" && head.ProtoName != "pb.LoginRoomRequest" {
		return nil
	}

	room, err := NewRoom(self, head.RoomId)
	if err != nil {
		self.Log.Error("big err create room err %s", head.RoomId)
		return nil
	}

	self.roomMutex.Lock()
	checkRoom, ok := self.id2Room[head.RoomId]
	if ok {
		return checkRoom
	}

	self.id2Room[head.RoomId] = room
	self.roomMutex.Unlock()

	go room.Run()
	return room
}

func (self *RoomMgr) MakeRoomEnd(roomId string) {
	self.roomMutex.Lock()
	delete(self.id2Room, roomId)
	self.roomMutex.Unlock()
}
