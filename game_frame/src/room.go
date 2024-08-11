package game_frame

import (
	"fmt"
	"game_frame/src/constants"
	"game_frame/src/domain"
	"game_frame/src/internal"
	"game_frame/src/pb"
	"github.com/golang/protobuf/proto"
	"math/rand"
	"reflect"
	"time"
)

const (
	ROOM_FPS            = 8
	ROOM_FRAME_INTERVAL = time.Second / ROOM_FPS

	ROOM_WAITING_PLAYERS_READY_TIME = time.Second * 20
)

type Room struct {
	RoomMgr *RoomMgr
	Counter *AtomicCounter

	GameInfo       *domain.GameInfo
	FrameId        int
	GameRunFrameId int //游戏运行的帧数

	uid2PlayerInfo  map[string]*domain.GamePlayer
	Players         []*domain.GamePlayer
	frameTicker     *time.Ticker
	MsgFromMgr      chan Closure
	rand            *rand.Rand
	protocol2Method map[string]reflect.Value

	State             string
	RoomId            string
	StateTimeoutFrame int
	ForceExitFrame    int
	isHall            bool
	IsAi              bool
	isInit            bool
	AiUid             string //AI的UID
	WinUid            string //赢的UID

	exit chan bool
}

var validRoomProtocols = map[string]string{
	"pb.LoginHallRequest":    "LoginHallHandler",
	"pb.LoadProgressRequest": "LoadProgressHandler",
	"pb.PingRequest":         "PingHandler",

	// 添加Room允许client访问的成员函数名
}

func NewRoom(roomMgr *RoomMgr, RoomId string) (*Room, error) {
	timer := time.NewTimer(0)
	timer.Stop()

	now := time.Now()
	self := &Room{
		RoomMgr:  roomMgr,
		GameInfo: roomMgr.Server.DynamicConfig.GameInfo,

		FrameId:           0,
		GameRunFrameId:    0,
		StateTimeoutFrame: 0,
		ForceExitFrame:    0,

		// 登录成功，则会重置为0；登录失败则一段时间后自动结束房间。+5s是为了防止WaitReconnectDuration配置成0时，创建房间就自动结束了
		Players:     make([]*domain.GamePlayer, 0),
		frameTicker: nil,
		rand:        rand.New(rand.NewSource(now.Unix())),
		State:       constants.ROOM_STATE_LOAD,
		RoomId:      RoomId,
		AiUid:       "",
		WinUid:      "",

		protocol2Method: make(map[string]reflect.Value),

		MsgFromMgr: make(chan Closure, 1024*10),

		Counter: &AtomicCounter{},

		IsAi:   false,
		isInit: false,
		isHall: false,

		exit: make(chan bool, 1),
	}

	t := reflect.ValueOf(self)

	protocol2Method := make(map[string]reflect.Value)
	for p, v := range validRoomProtocols {
		m := t.MethodByName(v)
		if !m.IsValid() {
			internal.GLog.Error("error, protocol handler not found, %s", p)
		}
		protocol2Method[p] = m
	}
	self.protocol2Method = protocol2Method
	internal.GLog.Info("NewRoom protocol2Method %+v", self.protocol2Method)

	return self, nil
}

func (self *Room) ApplyProtoHandler(reply string, head *pb.CommonHead, data []byte) {
	internal.GLog.Info("ApplyProtoHandler [%+v] ", head.PbName)
	method, ok := self.protocol2Method[head.PbName]
	if !ok {
		internal.GLog.Info("method not found: %+v", head.PbName)
		commonRes := &pb.GameCommonResponse{}
		commonRes.Code = constants.SYSTEM_ERROR
		commonRes.Msg = "system err"
		res, _ := proto.Marshal(commonRes)
		internal.NatsPool.Publish(reply, map[string]interface{}{"res": "ok", "data": string(res)})

		return
	}

	in := []reflect.Value{
		reflect.ValueOf(reply),
		reflect.ValueOf(head),
		reflect.ValueOf(data),
	}
	method.Call(in)
}

func (self *Room) Run() {
	defer func() {
		p := recover()
		if p != nil {
			internal.GLog.Info(fmt.Sprintf("execute panic recovered and going to stop: %v", p))
		}
	}()

	self.frameTicker = time.NewTicker(ROOM_FRAME_INTERVAL)
	self.StateTimeoutFrame = self.FrameId + int(ROOM_WAITING_PLAYERS_READY_TIME/ROOM_FRAME_INTERVAL)

	self.State = constants.ROOM_STATE_LOAD

	defer func() {
		self.frameTicker.Stop()
	}()

	for {
		// 优先查看exit，
		select {
		case <-self.exit:
			return
		case <-self.frameTicker.C:
			SafeRunClosure(self, func() {
				self.FrameId++

				if self.GameRunFrameId != 0 {
					self.GameRunFrameId++
				}
				self.Frame()
			})
		case c, _ := <-self.MsgFromMgr:
			SafeRunClosure(self, c)
		default:
			// do nothing
		}
	}
}

func (self *Room) Frame() {

	if self.checkReadyTimeOut() {
		internal.GLog.Info("roomId %v checkReadyTimeOut", self.RoomId)
		gameResult := self.getGameResult(constants.WIN_TYPE_NOT_START)
		self.enterEndState(gameResult, constants.WIN_TYPE_NOT_START)
		return
	}

	if self.checkTotalUseTimeOver() {
		internal.GLog.Info("roomId %v checkTotalUseTimeOver", self.RoomId)
		gameResult := self.getGameResult(constants.GAME_END_TOTAL_TIME_OVER)
		self.enterEndState(gameResult, constants.GAME_END_TOTAL_TIME_OVER)
		return
	}

	//判断是否断线超时
	if self.checkOffline() {
		gameResult := self.getGameResult(constants.GAME_END_OFFLINE)
		self.enterEndState(gameResult, constants.GAME_END_OFFLINE)
		return
	}

}

func (self *Room) SetState(state string) {
	self.State = state
}

func (self *Room) enterEndState(gameResult *pb.GameResult, reason string) {
	self.State = constants.ROOM_STATE_START
	//通知双方游戏结束

	internal.GLog.Info("roomId %v end reason %v ", self.RoomId, reason)
	self.onRoomEnd()
}

func (self *Room) onRoomEnd() {
	internal.GLog.Info("roomId %s end", self.RoomId)
	for _, playerInfo := range self.uid2PlayerInfo {
		if playerInfo != nil {
			//playerInfo.FlushCachedMsgs()
		}
	}

	self.RoomMgrRun(func(roomMgr *RoomMgr) {
		roomMgr.MakeRoomEnd(self.RoomId)
	})

	self.exit <- true
}

func (self *Room) RoomMgrRun(cb func(roomMgr *RoomMgr)) {
	roomMgr := self.RoomMgr
	roomMgr.MsgFromRoom <- func() {
		cb(roomMgr)
	}
}

func (self *Room) getGameResult(windType string) *pb.GameResult {
	return &pb.GameResult{}
}

func (self *Room) checkOffline() bool {
	if self.GameRunFrameId == 0 {
		return false
	}

	for _, player := range self.uid2PlayerInfo {
		if player.IsAi {
			continue
		}

		if player.OffLineFrameId == 0 {
			continue
		}

		if (self.GameRunFrameId-player.OffLineFrameId)/ROOM_FPS > 20 {
			internal.GLog.Info("roomId %v offline end uid %v ", self.RoomId, player.Uid)
			self.WinUid = self.getRivalUid(player.Uid)
			return true
		}
	}

	return false
}

func (self *Room) checkReadyTimeOut() bool {
	if self.State != constants.ROOM_STATE_LOAD {
		return false
	}

	if self.FrameId > self.StateTimeoutFrame {
		return true
	}

	return false
}

func (self *Room) checkTotalUseTimeOver() bool {
	if self.GameRunFrameId == 0 {
		return false
	}

	if int32(self.GameRunFrameId/ROOM_FPS) >= self.GameInfo.GameTime {
		self.WinUid = ""
		return true
	}

	return false
}

func (self *Room) SaveAndQuit() {
	self.exit <- true
}

func (self *Room) ResponseGateway(reply string, head *pb.CommonHead, response proto.Message) {
	head.Timestamp = time.Now().UnixMilli()
	bytes, _ := proto.Marshal(response)
	res := &pb.GameCommonResponse{
		Code: constants.CODE_SUCCESS,
		Msg:  "",
		Head: head,
		Data: string(bytes),
	}

	internal.GLog.Info("ResponseGateway reply: %+v protoName %+v", reply, head.GetPbName())

	commBytes, _ := proto.Marshal(res)
	internal.NatsPool.Publish(reply, map[string]interface{}{"res": "ok", "data": string(commBytes)})
}

// 房间内广播处理
func (self *Room) RoomBroadcast(msg proto.Message) {
	for _, player := range self.Players {
		self.Send2PlayerMessage(player, msg)
	}
}
