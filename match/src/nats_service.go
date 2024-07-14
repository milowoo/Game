package match

import (
	"github.com/golang/protobuf/proto"
	"github.com/nats-io/go-nats"
	"match/src/log"
	"match/src/mq"
	"match/src/pb"
)

type NatsService struct {
	Server             *Server
	log                *log.Logger
	Config             *GlobalConfig
	NatsPool           *mq.NatsPool
	GameMatchMgr       map[string]*GameMatch
	gameChange         chan string
	matchData          chan []byte
	cancelMatchData    chan []byte
	matchSubject       string
	cancelMatchSubject string
	exit               chan bool
}

func NewNatsService(server *Server) *NatsService {
	return &NatsService{
		Server:             server,
		log:                server.Log,
		Config:             server.Config,
		NatsPool:           server.NatsPool,
		matchData:          make(chan []byte),
		cancelMatchData:    make(chan []byte),
		GameMatchMgr:       make(map[string]*GameMatch),
		matchSubject:       "match.req",
		cancelMatchSubject: "match.cancel.req",
	}
}

func (self *NatsService) Run() {
	defer func() {
		p := recover()
		if p != nil {
			self.log.Info("execute panic recovered and going to stop: %v", p)
		}
	}()

	//监听消费匹配请求
	self.subscribeMatchRequest()
	//监听取消匹配请求
	self.subscribeCancelMatchRequest()

	for {
		// 优先查看exit，
		select {
		case <-self.exit:
			{
				return
			}
		case <-self.gameChange:
			{
				gameId := <-self.gameChange
				if self.GameMatchMgr[gameId] == nil {

				}
			}

		case <-self.matchData:
			{
				matchInfo := <-self.matchData
				self.ProcessMatchRequest(matchInfo)
			}
		case <-self.cancelMatchData:
			{
				cancelMatchReq := <-self.cancelMatchData
				self.ProcessCancelMatchRequest(cancelMatchReq)
			}

		default:
			// do nothing
		}
	}
}

func (self *NatsService) subscribeMatchRequest() {
	err := self.NatsPool.QueueSubscribe(self.matchSubject, "match_group", func(msg *nats.Msg) {
		self.matchData <- msg.Data
	})
	if err != nil {
		self.log.Error("subscribeMatchRequest err %+v", err)
	}
}

func (self *NatsService) subscribeCancelMatchRequest() {
	err := self.NatsPool.QueueSubscribe(self.cancelMatchSubject, "match_group", func(msg *nats.Msg) {
		self.cancelMatchData <- msg.Data
	})
	if err != nil {
		self.log.Error("subscribeCancelMatchRequest err %+v", err)
	}
}

func (self *NatsService) ProcessMatchRequest(data []byte) {
	var req pb.MatchRequest
	_ = proto.Unmarshal(data, &req)
	gameInfo := self.Server.DynamicConfig.GetGameInfo(req.GetGameId())
	if gameInfo == nil {
		self.log.Error("procMatchReq gameId %+v not found", req.GetGameId())
		return
	}

	gameMatch := self.GameMatchMgr[req.GetGameId()]
	if gameMatch == nil {
		self.log.Error("procMatchReq gameId %+v not found", req.GetGameId())
		return
	}

	RunOnMatch(gameMatch.MsgFromNats, gameMatch, func(gameMatch *GameMatch) {
		gameMatch.AddMatchRequest(&req)
	})
}

func (self *NatsService) ProcessCancelMatchRequest(data []byte) {
	var req pb.CancelMatchRequest
	_ = proto.Unmarshal(data, &req)
	gameInfo := self.Server.DynamicConfig.GetGameInfo(req.GetGameId())
	if gameInfo == nil {
		self.log.Error("ProcessCancelMatchRequest gameId %+v not found", req.GetGameId())
		return
	}

	gameMatch := self.GameMatchMgr[req.GetGameId()]
	if gameMatch == nil {
		self.log.Error("procMatchReq gameId %+v not found", req.GetGameId())
		return
	}

	RunOnMatch(gameMatch.MsgFromNats, gameMatch, func(gameMatch *GameMatch) {
		gameMatch.CancelMatchRequest(&req)
	})
}

func (self *NatsService) Quit() {
	self.exit <- true
}
