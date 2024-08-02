package gateway

import (
	"gateway/src/constants"
	"gateway/src/internal"
	"gateway/src/pb"
	"github.com/golang/protobuf/proto"
	"github.com/nats-io/go-nats"
)

type NatsGame struct {
	Server         *Server
	receiveSubject string
	exit           chan bool
}

func NewNatsGame(sever *Server) *NatsGame {
	return &NatsGame{
		Server:         sever,
		receiveSubject: constants.GetGamePushDataSubject(GetHostIp()),
		exit:           make(chan bool, 1),
	}
}

func (self *NatsGame) Run() {
	defer func() {
		p := recover()
		if p != nil {
			internal.GLog.Info("execute panic recovered and going to stop: %v", p)
		}
	}()

	internal.GLog.Info("nats game  begin ....")

	self.Server.WaitGroup.Add(1)
	defer func() {
		self.Server.WaitGroup.Done()
	}()

	for {
		select {
		case <-self.exit:
			{
				return
			}
		}
	}

	self.Quit()
}

func (self *NatsGame) Quit() {
	internal.GLog.Info("NatsGame quit")
	internal.NatsPool.Unsubscribe(self.receiveSubject)
	self.exit <- true
}

func (self *NatsGame) SubjectGamePushData() {
	internal.NatsPool.Subscribe(self.receiveSubject, func(mess *nats.Msg) {
		var pushData pb.GamePushMessage
		_ = proto.Unmarshal(mess.Data, &pushData)
		head := pushData.GetHead()
		internal.GLog.Info("SubjectGamePushData gameId %+v uid %+v pid %+v protoName %+v",
			head.GetGameId(), head.GetUid(), head.GetPid(), head.GetPbName())
		// uid 找出对应的 agent 进行匹配结果处理
		self.Server.AgentMgr.GamePushDataNotice(&pushData)
	})
}
