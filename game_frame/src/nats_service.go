package game_frame

import (
	"game_frame/src/constants"
	"game_frame/src/internal"
	"github.com/nats-io/go-nats"
)

type NatsService struct {
	Server                 *Server
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
			internal.GLog.Info("execute panic recovered and going to stop: %v", p)
		}
	}()

	self.Server.WaitGroup.Add(1)
	defer func() {
		self.Server.WaitGroup.Done()
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
	internal.GLog.Info("subscribeMatch subject %+v ", self.matchCreateRoomSubject)
	err := internal.NatsPool.SubscribeForRequest(self.matchCreateRoomSubject, func(subj, reply string, msg interface{}) {
		internal.GLog.Info("Nats subscribeMatch request subject:%+v,receive massage:%+v,reply subject:%+v", subj, msg, reply)

		roomId := msg.(string)
		hostIp := GetHostIp()

		internal.GLog.Info("subscribeMatch roomId %+v hostIp %+v", roomId, hostIp)
		internal.NatsPool.Publish(reply, map[string]interface{}{"res": "ok", "data": hostIp})
	})

	if err != nil {
		internal.GLog.Error("subscribeMatch err %+v", err)
	}

}

func (self *NatsService) subscribeGateway() {
	// 订阅一个Nats Request 主题
	internal.GLog.Info("subscribeGateway subject %+v ", self.gameSubject)
	err := internal.NatsPool.SubscribeForRequest(self.gameSubject, func(subj, reply string, msg interface{}) {
		internal.GLog.Info("subscribeGateway request subject:%+v,receive massage:%+v,reply subject:%+v", subj, msg, reply)

		roomMgr := self.Server.RoomMgr
		RunOnRoomMgr(roomMgr.MsgFromNats, roomMgr, func(roomMgr *RoomMgr) {
			roomMgr.ProcessGatewayRequest(reply, msg)
		})
	})

	if err != nil {
		internal.GLog.Error("subscribeGateway err %+v", err)
	}
}

func (self *NatsService) subscribeGM() {
	err := internal.NatsPool.Subscribe(self.gmSubject, func(msg *nats.Msg) {
		self.gmData <- msg.Data
	})
	if err != nil {
		internal.GLog.Error("subscribeMatch err %+v", err)
	}
}

func (self *NatsService) processGmCode(data []byte) {

}

func (self *NatsService) Quit() {
	internal.GLog.Info("nats service quit ....")
	internal.NatsPool.Unsubscribe(self.matchCreateRoomSubject, self.gameSubject, self.gmSubject)

	self.exit <- true
}
