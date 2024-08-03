package gateway

import (
	"gateway/src/constants"
	"gateway/src/internal"
	"gateway/src/pb"
	"github.com/golang/protobuf/proto"
	"github.com/nats-io/go-nats"
	"strconv"
	"time"
)

type NatsMatch struct {
	Server         *Server
	receiveSubject string
	MsgFromServer  chan Closure
	exit           chan bool
	isQuit         bool
}

func NewNatsMatch(sever *Server) *NatsMatch {
	return &NatsMatch{
		Server:         sever,
		MsgFromServer:  make(chan Closure, 2*1024),
		receiveSubject: constants.GetMatchResultSubject(GetHostIp()),
		exit:           make(chan bool, 1),
		isQuit:         false,
	}
}

func (self *NatsMatch) MatchRequest(gameId string, uid string, score string, opt string) error {
	matchReq := &pb.MatchRequest{
		GameId:         gameId,
		Uid:            uid,
		Score:          score,
		TimeStamp:      strconv.FormatInt(time.Now().UnixMilli(), 10),
		ReceiveSubject: self.receiveSubject,
		Opt:            opt,
	}

	request, _ := proto.Marshal(matchReq)
	var response interface{}
	err := internal.NatsPool.Request(constants.MATCH_SUBJECT, string(request), &response, 3*time.Second)
	if err != nil {
		internal.GLog.Error("MatchRequest gameId %+v uid %+v match err %+v", err)
		return err
	}
	return nil
}

func (self *NatsMatch) SubjectMatchResponse() {
	internal.GLog.Info("SubjectMatchResponse subject %+v", self.receiveSubject)
	internal.NatsPool.Subscribe(self.receiveSubject, func(mess *nats.Msg) {
		var matchOverRes pb.MatchOverRes
		_ = proto.Unmarshal(mess.Data, &matchOverRes)
		internal.GLog.Info("SubjectMatchResponse %+v", matchOverRes)
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
	var response interface{}
	err := internal.NatsPool.Request(constants.CANCEL_MATCH_SUBJECT, string(request), &response, 3*time.Second)
	if err != nil {
		internal.GLog.Error("MatchRequest gameId %+v uid %+v match err %+v", gameId, uid, err)
		return err
	}
	return nil
}

func (self *NatsMatch) Run() {
	defer func() {
		p := recover()
		if p != nil {
			internal.GLog.Info("execute panic recovered and going to stop: %v", p)
		}
	}()

	internal.GLog.Info("nats match  begin ....")

	self.SubjectMatchResponse()

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

	internal.GLog.Info("NatsMatch quit")
	self.isQuit = true
	internal.NatsPool.Unsubscribe(self.receiveSubject)
	self.exit <- true
}
