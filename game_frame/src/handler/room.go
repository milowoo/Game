package handler

import (
	"fmt"
	"game_frame/src"
	"game_frame/src/constants"
	"game_frame/src/log"
	"game_frame/src/pb"
	"game_frame/src/redis"
	"github.com/gogo/protobuf/proto"
	"math/rand"
	"reflect"
	"time"
)

const (
	ROOM_FPS            = 8
	ROOM_FRAME_INTERVAL = time.Second / ROOM_FPS

	ROOM_WAITING_PLAYERS_READY_TIME = time.Second * 20
)

type GamePlayer struct {
	Uid             string
	Pid             string
	RoomId          string
	HallId          string //到玩法服才有效
	LoadingProgress int32
	TotalUseTime    int
	Score           int
	IsAi            bool
	IsNewPlayer     bool
}

func NewPlayer(roomId string, uid string, pid string, isAi bool) *GamePlayer {
	p := &GamePlayer{
		Uid:    uid,
		Pid:    pid,
		RoomId: roomId,
		HallId: "",

		IsAi:            isAi,
		IsNewPlayer:     false,
		TotalUseTime:    0,
		Score:           0,
		LoadingProgress: 0,
	}

	return p
}

type Room struct {
	GlobalConfig  *game_frame.GlobalConfig
	DynamicConfig *game_frame.DynamicConfig

	RedisDao *redis.RedisDao

	RoomMgr *game_frame.RoomMgr
	Log     *log.Logger

	frameId        int
	gameRunFrameId int //游戏运行的帧数
	offLineFrameId int //掉线的帧数

	uid2PlayerInfo  map[string]*GamePlayer
	Players         []*GamePlayer
	frameTicker     *time.Ticker
	MsgFromMgr      chan game_frame.Closure
	rand            *rand.Rand
	protocol2Method map[string]reflect.Value

	state              string
	RoomId             string
	lastHeartBeatFrame int
	stateTimeoutFrame  int
	forceExitFrame     int
	TotalPoints        int
	isHall             bool
	IsAi               bool
	isInit             bool
	AiUid              string //AI的UID
	WinUid             string //赢的UID

	exit chan bool
}

var validRoomProtocols = map[string]string{
	"pb.c2sLoading":       "C2SLoading",
	"pb.LoginHallRequest": "LoginHall",

	// 添加Room允许client访问的成员函数名
}

func NewRoom(roomMgr *game_frame.RoomMgr, RoomId string) (*Room, error) {
	timer := time.NewTimer(0)
	timer.Stop()

	now := time.Now()
	self := &Room{
		RoomMgr:       roomMgr,
		RedisDao:      roomMgr.RedisDao,
		GlobalConfig:  roomMgr.GlobalConfig,
		DynamicConfig: roomMgr.DynamicConfig,

		frameId:            0,
		gameRunFrameId:     0,
		stateTimeoutFrame:  0,
		lastHeartBeatFrame: 0,
		forceExitFrame:     0,

		// 登录成功，则会重置为0；登录失败则一段时间后自动结束房间。+5s是为了防止WaitReconnectDuration配置成0时，创建房间就自动结束了
		Players:     make([]*GamePlayer, 0),
		frameTicker: nil,
		rand:        rand.New(rand.NewSource(now.Unix())),
		state:       constants.ROOM_STATE_LOAD,
		RoomId:      RoomId,
		AiUid:       "",
		WinUid:      "",

		MsgFromMgr: make(chan game_frame.Closure, 1024*10),

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
			self.Log.Error("error, protocol handler not found, %s", p)
		}
		protocol2Method[p] = m
	}
	self.protocol2Method = protocol2Method

	return self, nil
}

func (self *Room) ApplyProtoHandler(reply string, head *pb.CommonHead, msg proto.Message) {
	method, ok := self.protocol2Method[head.ProtoName]
	if !ok {
		self.Log.Info("method not found: %s", head.ProtoName)
		return
	}

	in := []reflect.Value{
		reflect.ValueOf(reply),
		reflect.ValueOf(head),
		reflect.ValueOf(msg),
	}
	method.Call(in)
}

func (self *Room) Run() {
	defer func() {
		p := recover()
		if p != nil {
			self.Log.Info(fmt.Sprintf("execute panic recovered and going to stop: %v", p))
		}
	}()

	self.frameTicker = time.NewTicker(ROOM_FRAME_INTERVAL)

	self.state = constants.ROOM_STATE_LOAD
	//self.stateTimeoutFrame = self.frameId + int32(ROOM_WAITING_READY_TIME/ROOM_FRAME_INTERVAL)

	defer func() {
		self.frameTicker.Stop()
	}()

	for {
		// 优先查看exit，
		select {
		case <-self.exit:
			return
		default:
			// do nothing
		}

		select {
		case <-self.frameTicker.C:
			game_frame.SafeRunClosure(self, func() {
				self.frameId++

				if self.gameRunFrameId != 0 {
					self.gameRunFrameId++
				}
				self.Frame()
			})
		}
	}
}

func (self *Room) Frame() {

}
func (self *Room) SaveAndQuit() {

}

func (self *Room) ResponseGateway(reply string, head *pb.CommonHead, response proto.Message) {
	head.Timestamp = time.Now().Unix()
	bytes, _ := proto.Marshal(response)
	res := &pb.GameCommonResponse{
		Code: constants.CODE_SUCCESS,
		Msg:  "",
		Head: head,
		Data: bytes,
	}

	commBytes, _ := proto.Marshal(res)
	self.RoomMgr.Server.NatsPool.Publish(reply, commBytes)
}

// 房间内广播处理
func (self *Room) RoomBroadcast(msg proto.Message) {
	//for _, player := range self.players {
	//	if player.IsAi {
	//		continue
	//	}
	//
	//	player.Ctx.SendInRoomPB(player.Uid, msg)
	//}
}
