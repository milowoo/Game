package game_frame

import (
	"game_frame/src/handler"
	"game_frame/src/log"
	"game_frame/src/pb"
	"game_frame/src/redis"
	"github.com/gogo/protobuf/proto"
	"github.com/nats-io/go-nats"
	"math/rand"
	"sync"
	"time"
)

const (
	RoomMgrFrameInterval = time.Millisecond * 500
	LogRoomCountInterval = time.Second * 30
	MatchTimeInterval    = time.Second * 5
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
	id2Room               map[string]*handler.Room
	innerId2Room          map[string]*handler.Room
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
		id2Room:    make(map[string]*handler.Room, 0),

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
	natsService := self.Server.NatsService
	for _, room := range self.id2Room {
		RunOnRoom(room.MsgFromMgr, room, func(input *handler.Room) {
			input.SaveAndQuit()
		})

		if natsService != nil {
			RunOnNatsService(natsService.MsgFromRoomMgr, natsService, func(natsService *NatsService) {
				natsService.UnsubscribeRoom(room.RoomId)
			})
		}
	}

	for len(self.id2Room) > 0 {
		c := <-self.MsgFromRoom
		SafeRunClosure(self, c)
	}
}

func (self *RoomMgr) ProcessGatewayRequest(reply string, msg *nats.Msg) {
	var commonReq pb.GameCommonRequest
	commonRes := &pb.GameCommonResponse{}
	proto.Unmarshal(msg.Data, &commonReq)
	if commonReq.GetGameId() != self.Server.Config.GameId {
		self.Log.Error("ProcessGatewayRequest invalid nat data %+v ", msg.Data)
		commonRes.Code = 101
		commonRes.Msg = "Unmarshal request err"
		commonRes.Timestamp = time.Now().Unix()
		res, _ := proto.Marshal(commonRes)
		self.Server.NatsPool.Publish(reply, res)
		return
	}

	commonRes.Sn = commonReq.Sn
	commonRes.ProtoName = commonReq.ProtoName

	var request proto.Message
	proto.Unmarshal(commonReq.GetData(), request)
	//根据房间ID 查找 room 信息
	room, ok := self.id2Room[commonReq.GetRoomId()]
	if ok {
		room.ApplyProtoHandler(commonRes, reply, commonReq.GetUid(), commonReq.GetProtoName(), request)
		return
	}

	room = self.GetOrCreateRoom(commonReq.GetRoomId(), commonReq.GetProtoName())
	if room == nil {
		self.Log.Error("ProcessGatewayRequest invalid proto %+v gameId %+v uid %+v",
			commonReq.GetProtoName(), commonReq.GetRoomId(), commonReq.GetUid())
		return
	}

	room.ApplyProtoHandler(commonRes, reply, commonReq.GetUid(), commonReq.GetProtoName(), request)
	return
}

func (self *RoomMgr) GetOrCreateRoom(roomId string, pbName string) *handler.Room {
	//先判断房间是否在使用中
	room, ok := self.id2Room[roomId]
	if ok {
		return room
	}

	//如果不是进入房间/ 大厅的协议，就是非法的请求
	if pbName != "pb.LoginHallRequest" && pbName != "pb.LoginRoomRequest" {
		return nil
	}

	room, err := handler.NewRoom(self, roomId)
	if err != nil {
		self.Log.Error("big err create room err %s", roomId)
		return nil
	}

	self.roomMutex.Lock()
	checkRoom, ok := self.id2Room[roomId]
	if ok {
		return checkRoom
	}

	self.id2Room[roomId] = room
	self.roomMutex.Unlock()

	go room.Run()
	return room
}

func (self *RoomMgr) MakeRoomEnd(roomId string) {
	self.roomMutex.Lock()
	delete(self.id2Room, roomId)
	self.roomMutex.Unlock()
	natsService := self.Server.NatsService
	RunOnNatsService(natsService.MsgFromRoomMgr, natsService, func(natsService *NatsService) {
		natsService.UnsubscribeRoom(roomId)
	})
}