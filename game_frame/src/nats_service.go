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
	gameSubject            string
	gmData                 chan []byte
	exit                   chan bool
}

func NewNatsService(server *Server) *NatsService {
	matchSubject := constants.GetCreateRoomNoticeSubject(server.Config.GameId, server.Config.GroupName)
	gmSubject := constants.GetGmCodeSubject(server.Config.GameId)
	gameSubject := constants.GetGameSubject(server.Config.GameId, GetHostIp())
	return &NatsService{
		Server:                 server,
		log:                    server.Log,
		Config:                 server.Config,
		NatsPool:               server.NatsPool,
		GameId:                 server.Config.GameId,
		MsgFromRoomMgr:         make(chan Closure, 2*1024),
		gmData:                 make(chan []byte, 128),
		matchCreateRoomSubject: matchSubject,
		gmSubject:              gmSubject,
		gameSubject:            gameSubject,
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

	self.subscribeGateway()

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

		default:
			// do nothing
		}
	}
}

func (self *NatsService) subscribeMatch() {
	// 订阅一个Nats Request 主题
	err := self.NatsPool.SubscribeForRequest(self.matchCreateRoomSubject, func(subj, reply string, msg interface{}) {
		self.log.Info("Nats subscribeMatch request subject:%+v,receive massage:%+v,reply subject:%+v", subj, msg, reply)

		natsMsg, ok := msg.(*nats.Msg)
		if ok {
			roomId := string(natsMsg.Data)
			hostIp := GetHostIp()

			self.log.Error("subscribeMatch roomId %+v hostIp %+v", roomId, hostIp)
			self.NatsPool.Publish(reply, []byte(hostIp))
		}

	})

	if err != nil {
		self.log.Error("subscribeMatch err %+v", err)
	}

}

func (self *NatsService) subscribeGateway() {
	// 订阅一个Nats Request 主题
	err := self.NatsPool.SubscribeForRequest(self.gameSubject, func(subj, reply string, msg interface{}) {
		self.log.Info("Nats Subscribe request subject:%+v,receive massage:%+v,reply subject:%+v", subj, msg, reply)

		natsMsg, ok := msg.(*nats.Msg)
		if ok {
			roomMgr := self.Server.RoomMgr
			RunOnRoomMgr(roomMgr.MsgFromNats, roomMgr, func(roomMgr *RoomMgr) {
				roomMgr.ProcessGatewayRequest(reply, natsMsg)
			})
			self.log.Error("subscribeGateway Failed to convert interface{} to *nats.Msg")
		}

	})

	if err != nil {
		self.log.Error("subscribeGateway err %+v", err)
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
}

func (self *NatsService) Quit() {
	self.NatsPool.Unsubscribe(self.matchCreateRoomSubject)
	self.NatsPool.Unsubscribe(self.gmSubject)

	self.exit <- true
}
