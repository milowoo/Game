package match

import (
	"github.com/golang/protobuf/proto"
	"match/src/constants"
	"match/src/log"
	"match/src/mq"
	"match/src/pb"
)

type NatsService struct {
	Server   *Server
	log      *log.Logger
	Config   *GlobalConfig
	NatsPool *mq.NatsPool
	exit     chan bool
}

func NewNatsService(server *Server) *NatsService {
	return &NatsService{
		Server:   server,
		log:      server.Log,
		Config:   server.Config,
		NatsPool: server.NatsPool,
		exit:     make(chan bool, 1),
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

	self.Server.WaitGroup.Add(1)
	defer func() {
		self.Server.WaitGroup.Done()
	}()

	for {
		// 优先查看exit，
		select {
		case <-self.exit:
			{
				return
			}

		default:
			// do nothing
		}
	}
	self.log.Info("nats service exit")
}

func (self *NatsService) subscribeLoginHallRequest() {
	// 订阅Nats 进入大厅 主题
	err := self.NatsPool.SubscribeForRequest(constants.LOGIN_HALL_SUBJECT, func(subj, reply string, msg interface{}) {
		self.log.Info("subscribeLoginHallRequest Subscribe request subject:%+v,receive massage:%+v,reply subject:%+v", subj, msg, reply)
		req, _ := ConvertInterfaceToString(msg)
		var request pb.CreateHallRequest
		proto.Unmarshal([]byte(req), &request)
		gameMatch := self.Server.MatchMgr[request.GetGameId()]
		self.log.Info("subscribeLoginHallRequest adasd")
		if gameMatch == nil {
			self.log.Error("subscribeLoginHallRequest gameId %+v not found", request.GetGameId())
			response := &pb.CreateHallResponse{
				Code: constants.INVALID_GAME_ID,
				Msg:  "invalid gameId",
			}

			res, _ := proto.Marshal(response)
			self.Server.NatsPool.Publish(reply, map[string]interface{}{"res": "ok", "data": string(res)})
			return
		}

		self.log.Info("subscribeLoginHallRequest 11111")

		RunOnMatch(gameMatch.MsgFromNats, gameMatch, func(gameMatch *GameMatch) {
			self.log.Info("subscribeLoginHallRequest 22222")
			gameMatch.LoginHallRequest(reply, &request)
		})
	})
	if err != nil {
		self.log.Error("subscribeLoginHallRequest err %+v", err)
	}
}

func (self *NatsService) subscribeMatchRequest() {
	// 订阅Nats 匹配 主题
	err := self.NatsPool.SubscribeForRequest(constants.MATCH_SUBJECT, func(subj, reply string, msg interface{}) {
		self.log.Info("subscribeMatchRequest Subscribe request subject:%+v,receive massage:%+v,reply subject:%+v", subj, msg, reply)
		req, err := ConvertInterfaceToString(msg)
		if err == nil {
			var request pb.MatchRequest
			proto.Unmarshal([]byte(req), &request)
			gameMatch := self.Server.MatchMgr[request.GetGameId()]
			if gameMatch == nil {
				self.log.Error("subscribeMatchRequest gameId %+v not found", request.GetGameId())
				response := &pb.MatchResponse{
					Code: constants.INVALID_GAME_ID,
					Msg:  "invalid gameId",
				}

				res, _ := proto.Marshal(response)
				self.Server.NatsPool.Publish(reply, map[string]interface{}{"res": "ok", "data": string(res)})
				return
			}

			RunOnMatch(gameMatch.MsgFromNats, gameMatch, func(gameMatch *GameMatch) {
				gameMatch.AddMatchRequest(&request)
				response := &pb.MatchResponse{
					Code: constants.CODE_SUCCESS,
				}

				res, _ := proto.Marshal(response)
				self.Server.NatsPool.Publish(reply, map[string]interface{}{"res": "ok", "data": string(res)})
			})
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
		req, err := ConvertInterfaceToString(msg)
		if err == nil {
			var request pb.CancelMatchRequest
			proto.Unmarshal([]byte(req), &request)
			gameMatch := self.Server.MatchMgr[request.GetGameId()]
			if gameMatch == nil {
				self.log.Error("subscribeCancelMatchRequest gameId %+v not found", request.GetGameId())
				response := &pb.CancelMatchResponse{
					Code: constants.INVALID_GAME_ID,
					Msg:  "invalid gameId",
				}

				res, _ := proto.Marshal(response)
				self.Server.NatsPool.Publish(reply, map[string]interface{}{"res": "ok", "data": string(res)})
				return
			}

			RunOnMatch(gameMatch.MsgFromNats, gameMatch, func(gameMatch *GameMatch) {
				gameMatch.CancelMatchRequest(&request)
				response := &pb.CancelMatchResponse{
					Code: constants.CODE_SUCCESS,
				}
				res, _ := proto.Marshal(response)
				self.Server.NatsPool.Publish(reply, map[string]interface{}{"res": "ok", "data": string(res)})
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
	self.log.Info("Quit NatsService ....")
	self.exit <- true
}
