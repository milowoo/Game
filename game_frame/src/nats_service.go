package game_frame

import (
	"game_frame/src/constants"
	"game_frame/src/log"
	"game_frame/src/mq"
	"github.com/nats-io/go-nats"
)

type NatsService struct {
	Server                 *Server
	log                    *log.Logger
	Config                 *GlobalConfig
	NatsPool               *mq.NatsPool
	MsgFromRoomMgr         chan Closure
	GameId                 string
	matchCreateRoomSubject string
	gmSubject              string
	roomSubjectMap         map[string]int32
	matchData              chan []byte
	gmData                 chan []byte
	exit                   chan bool
}

func NewNatsService(server *Server) *NatsService {
	matchSubject := constants.GetCreateRoomNoticeSubject(server.Config.GameId, server.Config.GroupName)
	gmSubject := constants.GetGmCodeSubject(server.Config.GameId)
	return &NatsService{
		Server:                 server,
		log:                    server.Log,
		Config:                 server.Config,
		NatsPool:               server.NatsPool,
		GameId:                 server.Config.GameId,
		MsgFromRoomMgr:         make(chan Closure, 2*1024),
		matchData:              make(chan []byte, 1024),
		gmData:                 make(chan []byte, 128),
		roomSubjectMap:         make(map[string]int32),
		matchCreateRoomSubject: matchSubject,
		gmSubject:              gmSubject,
		exit:                   make(chan bool, 1),
	}
}

func (self *NatsService) Run() {
	defer func() {
		p := recover()
		if p != nil {
			self.log.Info("execute panic recovered and going to stop: %v", p)
		}
	}()

	//监听匹配创建房间信息
	self.subscribeMatch()
	//监听GM信息
	self.subscribeGM()

	for {
		// 优先查看exit，
		select {
		case <-self.exit:
			{
				return
			}
		case c := <-self.MsgFromRoomMgr:
			SafeRunClosure(self, c)
		case <-self.gmData:
			{
				gmData := <-self.gmData
				self.processGmCode(gmData)
			}

		case <-self.matchData:
			{
				roomId := <-self.matchData
				self.ProcessCreateRoom(string(roomId))
			}

		default:
			// do nothing
		}
	}
}

func (self *NatsService) subscribeMatch() {
	err := self.NatsPool.QueueSubscribe(self.matchCreateRoomSubject, "game_match_group", func(msg *nats.Msg) {
		self.matchData <- msg.Data
	})
	if err != nil {
		self.log.Error("subscribeMatch err %+v", err)
	}
}

func (self *NatsService) ProcessCreateRoom(roomId string) {
	//对房间进行监听， 处理gateway的消息 "game." + server.Config.GameId + ".*",
	subject := constants.GetGameSubject(self.GameId, roomId)
	self.roomSubjectMap[roomId] = 1
	//监听gateway请求数据
	self.subscribeGateway(subject)
}

func (self *NatsService) subscribeGateway(subject string) {
	// 订阅一个Nats Request 主题
	err := self.NatsPool.SubscribeForRequest(subject, func(subj, reply string, msg interface{}) {
		self.log.Info("Nats Subscribe request subject:%+v,receive massage:%+v,reply subject:%+v", subj, msg, reply)

		natsMsg, ok := msg.(*nats.Msg)
		if ok {
			roomMgr := self.Server.RoomMgr
			RunOnRoomMgr(roomMgr.MsgFromNats, roomMgr, func(roomMgr *RoomMgr) {
				roomMgr.ProcessGatewayRequest(reply, natsMsg)
			})
			self.log.Error("SubscribeGetUid Failed to convert interface{} to *nats.Msg")
		}

	})

	if err != nil {
		self.log.Error("SubscribeGetUid err %+v", err)
	}
}

func (self *NatsService) subscribeGM() {
	err := self.NatsPool.Subscribe(self.gmSubject, func(msg *nats.Msg) {
		self.gmData <- msg.Data
	})
	if err != nil {
		self.log.Error("subscribeMatch err %+v", err)
	}
}

func (self *NatsService) processGmCode(data []byte) {

}

func (self *NatsService) getRoomSubject(roomId string) string {
	subject := "game." + self.Config.GameId + "." + roomId
	return subject
}

func (self *NatsService) UnsubscribeRoom(roomId string) {
	subject := self.getRoomSubject(roomId)
	self.NatsPool.Unsubscribe(subject)
	delete(self.roomSubjectMap, roomId)
}

func (self *NatsService) Quit() {
	self.NatsPool.Unsubscribe(self.matchCreateRoomSubject)
	self.NatsPool.Unsubscribe(self.gmSubject)
	for roomId, _ := range self.roomSubjectMap {
		subject := self.getRoomSubject(roomId)
		self.NatsPool.Unsubscribe(subject)
	}
	self.exit <- true
}
