package gateway

import (
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
	matchSubject   string
	cancelSubject  string
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
		matchSubject:   "match.req",
		cancelSubject:  "match.cancel.req",
		NatsPool:       sever.NatsPool,
		MsgFromServer:  make(chan Closure, 2*1024),
		receiveSubject: "match.response." + GetHostIp(),
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

	err := self.NatsPool.Publish(self.matchSubject, request)
	if err != nil {
		self.log.Error("MatchRequest gameId %+v uid %+v match err %+v", err)
		return err
	}
	return nil
}

func (self *NatsMatch) MatchResponse() {
	self.NatsPool.Subscribe(self.receiveSubject, func(mess *nats.Msg) {
		var matchOverRes pb.MatchOverRes
		_ = proto.Unmarshal(mess.Data, &matchOverRes)
		self.log.Info("MatchResponse %+v", matchOverRes)
		// uid 找出对应的 agent 进行匹配结果处理
		self.Server.AgentMgr.MatchResponse(&matchOverRes)
	})
}

func (self *NatsMatch) CancelMatchRequest(gameId string, uid string) error {
	matchReq := &pb.CancelMatchRequest{
		GameId: gameId,
		Uid:    uid,
	}

	request, _ := proto.Marshal(matchReq)
	err := self.NatsPool.Publish(self.cancelSubject, request)
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
