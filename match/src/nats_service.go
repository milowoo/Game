package match

import (
	"github.com/golang/protobuf/proto"
	"github.com/nats-io/go-nats"
	"match/src/constants"
	"match/src/log"
	"match/src/mq"
	"match/src/pb"
)

type NatsService struct {
	Server       *Server
	log          *log.Logger
	Config       *GlobalConfig
	NatsPool     *mq.NatsPool
	GameMatchMgr map[string]*GameMatch
	gameChange   chan string
	exit         chan bool
}

func NewNatsService(server *Server) *NatsService {
	return &NatsService{
		Server:       server,
		log:          server.Log,
		Config:       server.Config,
		NatsPool:     server.NatsPool,
		GameMatchMgr: make(map[string]*GameMatch),
		exit:         make(chan bool, 1),
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

	//监听消费进入大厅请求
	self.subscribeLoginHallRequest()

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

		default:
			// do nothing
		}
	}
}

func (self *NatsService) subscribeLoginHallRequest() {
	// 订阅Nats 进入大厅 主题
	err := self.NatsPool.SubscribeForRequest(constants.LOGIN_HALL_SUBJECT, func(subj, reply string, msg interface{}) {
		self.log.Info("Nats Subscribe request subject:%+v,receive massage:%+v,reply subject:%+v", subj, msg, reply)

		natsMsg, ok := msg.(*nats.Msg)
		if ok {
			var request pb.LoginHallRequest
			proto.Unmarshal(natsMsg.Data, &request)
			gameMatch := self.GameMatchMgr[request.GetGameId()]
			if gameMatch == nil {
				self.log.Error("subscribeLoginHallRequest gameId %+v not found", request.GetGameId())
				response := &pb.LoginHallResponse{
					Code: 102,
					Msg:  "invalid gameId",
				}

				res, _ := proto.Marshal(response)
				self.Server.NatsPool.Publish(reply, res)
				return
			}

			RunOnMatch(gameMatch.MsgFromNats, gameMatch, func(gameMatch *GameMatch) {
				gameMatch.LoginHallRequest(reply, &request)
			})
		} else {
			self.log.Error("subscribeLoginHallRequest Failed to convert interface{} to *nats.Msg")
		}

	})

	if err != nil {
		self.log.Error("subscribeLoginHallRequest err %+v", err)
	}
}

func (self *NatsService) subscribeMatchRequest() {
	// 订阅Nats 匹配 主题
	err := self.NatsPool.SubscribeForRequest(constants.MATCH_SUBJECT, func(subj, reply string, msg interface{}) {
		self.log.Info("subscribeMatchRequest Subscribe request subject:%+v,receive massage:%+v,reply subject:%+v", subj, msg, reply)

		natsMsg, ok := msg.(*nats.Msg)
		if ok {
			var request pb.MatchRequest
			proto.Unmarshal(natsMsg.Data, &request)
			gameMatch := self.GameMatchMgr[request.GetGameId()]
			if gameMatch == nil {
				self.log.Error("subscribeMatchRequest gameId %+v not found", request.GetGameId())
				response := &pb.MatchResponse{
					Code: 102,
					Msg:  "invalid gameId",
				}

				res, _ := proto.Marshal(response)
				self.Server.NatsPool.Publish(reply, res)
				return
			}

			RunOnMatch(gameMatch.MsgFromNats, gameMatch, func(gameMatch *GameMatch) {
				gameMatch.AddMatchRequest(&request)
			})
		} else {
			self.log.Error("subscribeMatchRequest Failed to convert interface{} to *nats.Msg")
		}

	})

	if err != nil {
		self.log.Error("subscribeMatchRequest err %+v", err)
	}

}

func (self *NatsService) subscribeCancelMatchRequest() {
	// 订阅Nats 取消匹配 主题
	err := self.NatsPool.SubscribeForRequest(constants.CANCEL_MATCH_SUBJECT, func(subj, reply string, msg interface{}) {
		self.log.Info("subscribeCancelMatchRequest Subscribe request subject:%+v,receive massage:%+v,reply subject:%+v", subj, msg, reply)

		natsMsg, ok := msg.(*nats.Msg)
		if ok {
			var request pb.CancelMatchRequest
			proto.Unmarshal(natsMsg.Data, &request)
			gameMatch := self.GameMatchMgr[request.GetGameId()]
			if gameMatch == nil {
				self.log.Error("subscribeCancelMatchRequest gameId %+v not found", request.GetGameId())
				response := &pb.CancelMatchResponse{
					Code: 102,
					Msg:  "invalid gameId",
				}

				res, _ := proto.Marshal(response)
				self.Server.NatsPool.Publish(reply, res)
				return
			}

			RunOnMatch(gameMatch.MsgFromNats, gameMatch, func(gameMatch *GameMatch) {
				gameMatch.CancelMatchRequest(&request)
			})
		} else {
			self.log.Error("subscribeCancelMatchRequest Failed to convert interface{} to *nats.Msg")
		}

	})

	if err != nil {
		self.log.Error("subscribeCancelMatchRequest err %+v", err)
	}
}

func (self *NatsService) Quit() {
	self.exit <- true
}
