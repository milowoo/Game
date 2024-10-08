package game_frame

import (
	"game_frame/src/constants"
	"game_frame/src/internal"
	"game_frame/src/pb"
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
	MsgFromNats chan Closure
	MsgFromRoom chan Closure

	rand      *rand.Rand
	roomMutex sync.Mutex

	frameTimer *time.Ticker
	id2Room    map[string]*Room
	frameId    int
	isQuit     bool
	exit       chan bool
}

func NewRoomMgr(sever *Server) *RoomMgr {
	now := time.Now()
	roomMgr := &RoomMgr{
		Server:      sever,
		MsgFromNats: make(chan Closure, 10*1024),
		MsgFromRoom: make(chan Closure, 2*1024),

		rand:       rand.New(rand.NewSource(now.Unix())),
		frameTimer: time.NewTicker(RoomMgrFrameInterval),
		id2Room:    make(map[string]*Room, 0),

		frameId: 0,
		isQuit:  false,
		exit:    make(chan bool, 1),
	}

	return roomMgr
}

func (self *RoomMgr) Run() {
	defer func() {
		p := recover()
		if p != nil {
			internal.GLog.Info("execute panic recovered and going to stop: %v", p)
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
				return
			}
		case c, ok := <-self.MsgFromNats:
			if !ok {
				self.Quit()
				return
			}
			SafeRunClosure(self, c)
		case c, _ := <-self.MsgFromRoom:
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
	if self.isQuit {
		return
	}

	self.isQuit = true

	// 通知所有room强制存储并退出
	for _, room := range self.id2Room {
		RunOnRoom(room.MsgFromMgr, room, func(input *Room) {
			input.SaveAndQuit()
		})
	}

	self.exit <- true
	internal.GLog.Info("room mgr quit ...")
}

func (self *RoomMgr) ProcessGatewayRequest(reply string, msg interface{}) {
	dataMap := msg.(map[string]interface{})
	head, data := ConvertRequest(dataMap)
	internal.GLog.Info("ProcessGatewayRequest uid %v  roomId %+v proto Name [%+v] hostIp [%v].....",
		head.GetUid(), head.GetRoomId(), head.GetPbName(), head.GetGatewayIp())

	if head.GetGameId() != self.Server.Config.GameId {
		commonRes := &pb.GameCommonResponse{}
		internal.GLog.Error("ProcessGatewayRequest invalid nat data %+v ", head.GameId)
		commonRes.Code = constants.INVALID_BODY
		commonRes.Msg = "Unmarshal request err"
		res, _ := proto.Marshal(commonRes)
		internal.NatsPool.Publish(reply, map[string]interface{}{"res": "ok", "data": string(res)})
		return
	}

	room := self.GetOrCreateRoom(head)
	if room == nil {
		commonRes := &pb.GameCommonResponse{}
		internal.GLog.Error("ProcessGatewayRequest invalid proto [%+v] roomId [%+v] uid [%+v]",
			head.GetPbName(), head.GetRoomId(), head.GetUid())
		commonRes.Code = constants.SYSTEM_ERROR
		commonRes.Msg = "system err"
		res, _ := proto.Marshal(commonRes)
		internal.NatsPool.Publish(reply, map[string]interface{}{"res": "ok", "data": string(res)})
		return
	}

	RunOnRoom(room.MsgFromMgr, room, func(input *Room) {
		room.ApplyProtoHandler(reply, head, data)
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
	if head.GetPbName() != "pb.LoginHallRequest" && head.GetPbName() != "pb.LoginRoomRequest" {
		internal.GLog.Info("GetOrCreateRoom invalid request proto name %+v", head.GetPbName())
		return nil
	}

	newRoom, err := NewRoom(self, head.RoomId)
	if err != nil {
		internal.GLog.Error("big err create room err %s", head.RoomId)
		return nil
	}

	self.roomMutex.Lock()
	checkRoom, ok := self.id2Room[head.RoomId]
	if ok {
		return checkRoom
	}

	self.id2Room[head.RoomId] = newRoom
	self.roomMutex.Unlock()

	go newRoom.Run()
	return newRoom
}

func (self *RoomMgr) MakeRoomEnd(roomId string) {
	self.roomMutex.Lock()
	delete(self.id2Room, roomId)
	self.roomMutex.Unlock()
}
