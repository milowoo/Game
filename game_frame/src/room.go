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

type GamePlayer struct {
	Uid             string
	Pid             string
	RoomId          string
	HallId          string
	GatewayIp       string
	LoadingProgress int32
	TotalUseTime    int
	Score           int
	IsAi            bool
	IsNewPlayer     bool
}

func NewPlayer(roomId string, uid string, pid string, hostIp string, isAi bool) *GamePlayer {
	p := &GamePlayer{
		Uid:       uid,
		Pid:       pid,
		RoomId:    roomId,
		HallId:    "",
		GatewayIp: hostIp,

		IsAi:            isAi,
		IsNewPlayer:     false,
		TotalUseTime:    0,
		Score:           0,
		LoadingProgress: 0,
	}

	return p
}

type Room struct {
	RoomMgr *RoomMgr
	Counter *AtomicCounter

	GameId         string
	GameInfo       *domain.GameInfo
	frameId        int
	gameRunFrameId int //游戏运行的帧数
	offLineFrameId int //掉线的帧数

	uid2PlayerInfo  map[string]*GamePlayer
	Players         []*GamePlayer
	frameTicker     *time.Ticker
	MsgFromMgr      chan Closure
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

func NewRoom(roomMgr *RoomMgr, RoomId string) (*Room, error) {
	timer := time.NewTimer(0)
	timer.Stop()

	now := time.Now()
	self := &Room{
		RoomMgr:  roomMgr,
		GameId:   roomMgr.Server.Config.GameId,
		GameInfo: roomMgr.Server.DynamicConfig.GameInfo,

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

	return self, nil
}

func (self *Room) ApplyProtoHandler(reply string, head *pb.CommonHead, data []byte) {
	internal.GLog.Info("ApplyProtoHandler reply: %s", reply)
	//method, ok := self.protocol2Method[head.ProtoName]
	//if !ok {
	//	self.Log.Info("method not found: %+v", head.ProtoName)
	//	commonRes := &pb.GameCommonResponse{}
	//	commonRes.Code = constants.SYSTEM_ERROR
	//	commonRes.Msg = "system err"
	//	res, _ := proto.Marshal(commonRes)
	//	self.NatPool.Publish(reply, map[string]interface{}{"res": "ok", "data": string(res)})
	//
	//	return
	//}

	if head.GetPbName() == "pb.LoginHallRequest" {
		self.LoginHall(reply, head, data)
	}

	//self.Log.Info("ApplyProtoHandler 1111")
	//in := []reflect.Value{
	//	reflect.ValueOf(reply),
	//	reflect.ValueOf(head),
	//	reflect.ValueOf(msg),
	//}
	//method.Call(in)
	//self.Log.Info("ApplyProtoHandler 2222")
}

func (self *Room) Run() {
	defer func() {
		p := recover()
		if p != nil {
			internal.GLog.Info(fmt.Sprintf("execute panic recovered and going to stop: %v", p))
		}
	}()

	self.frameTicker = time.NewTicker(ROOM_FRAME_INTERVAL)

	self.state = constants.ROOM_STATE_LOAD

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
				self.frameId++

				if self.gameRunFrameId != 0 {
					self.gameRunFrameId++
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
	self.exit <- true
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

func (self *Room) Send2PlayerMessage(player *GamePlayer, msg proto.Message) {
	if player.IsAi {
		return
	}

	typ := reflect.TypeOf(msg)
	protoName := typ.Elem().Name()
	bytes, _ := proto.Marshal(msg)

	head := &pb.PushHead{
		Uid:       player.Uid,
		GameId:    self.GameId,
		Pid:       player.Pid,
		Timestamp: time.Now().Unix(),
		RoomId:    self.RoomId,
		PbName:    protoName,
		Sn:        self.Counter.GetIncrementValue(),
	}

	data := &pb.GamePushMessage{
		Head: head,
		Data: string(bytes),
	}
	res, _ := proto.Marshal(data)
	internal.NatsPool.Publish(constants.GetGamePushDataSubject(player.GatewayIp), map[string]interface{}{"res": "ok", "data": string(res)})
}

func (self *Room) IsFirstLogin(uid string) bool {
	if len(self.Players) == 0 {
		return true
	}
	for _, player := range self.Players {
		if player.Uid == uid {
			return false
		}
	}

	return true
}

/*
*
获取竞争对手
*/
func (self *Room) GetRival(uid string) *GamePlayer {
	for _, player := range self.Players {
		if player.Uid != uid {
			return player
		}
	}

	return nil
}
