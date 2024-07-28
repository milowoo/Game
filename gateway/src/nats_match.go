package gateway

import (
	"gateway/src/constants"
	"gateway/src/log"
	"gateway/src/mq"
	"gateway/src/pb"
	"github.com/gogo/protobuf/proto"
	"github.com/nats-io/go-nats"
	"time"
)

type NatsMatch struct {
	Server         *Server
	log            *log.Logger
	receiveSubject string
	NatsPool       *mq.NatsPool
	MsgFromServer  chan Closure
	exit           chan bool
	isQuit         bool
}

func NewNatsMatch(sever *Server) *NatsMatch {
	return &NatsMatch{
		Server:         sever,
		log:            sever.Log,
		NatsPool:       sever.NatsPool,
		MsgFromServer:  make(chan Closure, 2*1024),
		receiveSubject: constants.GetMatchResultSubject(GetHostIp()),
		exit:           make(chan bool, 1),
		isQuit:         true,
	}
}

func (self *NatsMatch) MatchRequest(gameId string, uid string, score int32, opt string) error {
	matchReq := &pb.MatchRequest{
		GameId:         gameId,
		Uid:            uid,
		Score:          score,
		TimeStamp:      time.Now().Unix(),
		ReceiveSubject: self.receiveSubject,
		Opt:            opt,
	}

	request, _ := proto.Marshal(matchReq)

	err := self.NatsPool.Publish(constants.MATCH_SUBJECT, request)
	if err != nil {
		self.log.Error("MatchRequest gameId %+v uid %+v match err %+v", err)
		return err
	}
	return nil
}

func (self *NatsMatch) SubjectMatchResponse() {
	self.log.Info("SubjectMatchResponse subject %+v", self.receiveSubject)
	self.NatsPool.Subscribe(self.receiveSubject, func(mess *nats.Msg) {
		var matchOverRes pb.MatchOverRes
		_ = proto.Unmarshal(mess.Data, &matchOverRes)
		self.log.Info("SubjectMatchResponse %+v", matchOverRes)
		// uid 找出对应的 agent 进行匹配结果处理
		self.Server.AgentMgr.MatchResponse(&matchOverRes)
	})
}

func (self *NatsMatch) SubscribeGetUid() {
	self.log.Info("SubscribeGetUid subject %+v begin ... ", constants.UCENTER_APPLY_UID_SUBJECT)

	// 订阅一个Nats Request 主题
	err := self.NatsPool.SubscribeForRequest(constants.UCENTER_APPLY_UID_SUBJECT, func(subj, reply string, msg interface{}) {
		self.log.Info("Nats Subscribe request subject:%+v,receive massage:%+v,reply subject:%+v", subj, msg, reply)

		natsMsg, ok := msg.(*nats.Msg)
		if ok {
			self.log.Info("SubscribeGetUid 111 %+v", natsMsg.Subject)
		}
		//if ok {
		//	self.GetPlayerUID(reply, natsMsg)
		//} else {
		//	self.log.Error("SubscribeGetUid Failed to convert interface{} to *nats.Msg")
		//}

	})

	if err != nil {
		self.log.Error("SubscribeGetUid err %+v", err)
	}

}

func (self *NatsMatch) CancelMatchRequest(gameId string, uid string) error {
	matchReq := &pb.CancelMatchRequest{
		GameId: gameId,
		Uid:    uid,
	}

	request, _ := proto.Marshal(matchReq)
	err := self.NatsPool.Publish(constants.CANCEL_MATCH_SUBJECT, request)
	if err != nil {
		self.log.Error("MatchRequest gameId %+v uid %+v match err %+v", err)
		return err
	}
	return nil
}

func (self *NatsMatch) Run() {
	defer func() {
		p := recover()
		if p != nil {
			self.log.Info("execute panic recovered and going to stop: %v", p)
		}
	}()

	self.log.Info("nats match  begin ....")

	self.SubjectMatchResponse()

	self.SubscribeGetUid()

	self.Server.WaitGroup.Add(1)
	defer func() {
		self.Server.WaitGroup.Done()
	}()

ALL:
	for {
		select {
		case <-self.exit:
			{
				return
			}
		case c, ok := <-self.MsgFromServer:
			if !ok {
				break ALL
			}
			SafeRunClosure(self, c)
		}
	}

	self.Quit()
}

func (self *NatsMatch) Quit() {
	if self.isQuit {
		return
	}
	self.log.Info("NatsMatch quit")
	self.isQuit = true
	self.exit <- true
}
