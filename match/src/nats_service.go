package match

import (
	"github.com/golang/protobuf/proto"
	"match/src/config"
	"match/src/constants"
	"match/src/internal"
	"match/src/pb"
	"match/src/utils"
)

type NatsService struct {
	Server *Server
	Config *config.GlobalConfig
	exit   chan bool
}

func NewNatsService(server *Server) *NatsService {
	return &NatsService{
		Server: server,
		Config: server.Config,
		exit:   make(chan bool, 1),
	}
}

func (self *NatsService) Run() {
	defer func() {
		p := recover()
		if p != nil {
			internal.GLog.Info("execute panic recovered and going to stop: %v", p)
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
	internal.GLog.Info("nats service exit")
}

func (self *NatsService) subscribeLoginHallRequest() {
	// 订阅Nats 进入大厅 主题
	err := internal.NatsPool.SubscribeForRequest(constants.LOGIN_HALL_SUBJECT, func(subj, reply string, msg interface{}) {
		internal.GLog.Info("subscribeLoginHallRequest Subscribe request subject:%+v,receive massage:%+v,reply subject:%+v", subj, msg, reply)
		req, _ := utils.ConvertInterfaceToString(msg)
		var request pb.CreateHallRequest
		proto.Unmarshal([]byte(req), &request)
		gameMatch := self.Server.MatchMgr[request.GetGameId()]
		if gameMatch == nil {
			internal.GLog.Error("subscribeLoginHallRequest gameId %+v not found", request.GetGameId())
			response := &pb.CreateHallResponse{
				Code: constants.INVALID_GAME_ID,
				Msg:  "invalid gameId",
			}

			res, _ := proto.Marshal(response)
			internal.NatsPool.Publish(reply, map[string]interface{}{"res": "ok", "data": string(res)})
			return
		}

		RunOnMatch(gameMatch.MsgFromNats, gameMatch, func(gameMatch *GameMatch) {
			gameMatch.LoginHallRequest(reply, &request)
		})
	})
	if err != nil {
		internal.GLog.Error("subscribeLoginHallRequest err %+v", err)
	}
}

func (self *NatsService) subscribeMatchRequest() {
	// 订阅Nats 匹配 主题
	err := internal.NatsPool.SubscribeForRequest(constants.MATCH_SUBJECT, func(subj, reply string, msg interface{}) {
		var request pb.MatchRequest
		internal.GLog.Info("subscribeMatchRequest Subscribe request subject:%+v,receive massage:%+v,reply subject:%+v", subj, msg, reply)
		req, err := utils.ConvertInterfaceToString(msg)
		if err == nil {
			proto.Unmarshal([]byte(req), &request)
			internal.GLog.Info("subscribeMatchRequest %+v", &request)

			gameMatch := self.Server.MatchMgr[request.GetGameId()]
			if gameMatch == nil {
				internal.GLog.Error("subscribeMatchRequest gameId %+v not found", request.GetGameId())
				response := &pb.MatchResponse{
					Code: constants.INVALID_GAME_ID,
					Msg:  "invalid gameId",
				}

				res, _ := proto.Marshal(response)
				internal.NatsPool.Publish(reply, map[string]interface{}{"res": "ok", "data": string(res)})
				return
			}

			RunOnMatch(gameMatch.MsgFromNats, gameMatch, func(gameMatch *GameMatch) {
				response := &pb.MatchResponse{
					Code: constants.CODE_SUCCESS,
				}
				res, _ := proto.Marshal(response)
				internal.NatsPool.Publish(reply, map[string]interface{}{"res": "ok", "data": string(res)})
				gameMatch.AddMatchRequest(&request)
			})
		}
	})

	if err != nil {
		internal.GLog.Error("subscribeMatchRequest err %+v", err)
	}

}

func (self *NatsService) subscribeCancelMatchRequest() {
	// 订阅Nats 取消匹配 主题
	err := internal.NatsPool.SubscribeForRequest(constants.CANCEL_MATCH_SUBJECT, func(subj, reply string, msg interface{}) {
		internal.GLog.Info("subscribeCancelMatchRequest Subscribe request subject:%+v,receive massage:%+v,reply subject:%+v", subj, msg, reply)
		req, err := utils.ConvertInterfaceToString(msg)
		if err == nil {
			var request pb.CancelMatchRequest
			proto.Unmarshal([]byte(req), &request)
			gameMatch := self.Server.MatchMgr[request.GetGameId()]
			if gameMatch == nil {
				internal.GLog.Error("subscribeCancelMatchRequest gameId %+v not found", request.GetGameId())
				response := &pb.CancelMatchResponse{
					Code: constants.INVALID_GAME_ID,
					Msg:  "invalid gameId",
				}

				res, _ := proto.Marshal(response)
				internal.NatsPool.Publish(reply, map[string]interface{}{"res": "ok", "data": string(res)})
				return
			}

			RunOnMatch(gameMatch.MsgFromNats, gameMatch, func(gameMatch *GameMatch) {
				gameMatch.CancelMatchRequest(&request)
				response := &pb.CancelMatchResponse{
					Code: constants.CODE_SUCCESS,
				}
				res, _ := proto.Marshal(response)
				internal.NatsPool.Publish(reply, map[string]interface{}{"res": "ok", "data": string(res)})
			})
		} else {
			internal.GLog.Error("subscribeCancelMatchRequest Failed to convert interface{} to *nats.Msg")
		}

	})

	if err != nil {
		internal.GLog.Error("subscribeCancelMatchRequest err %+v", err)
	}
}

func (self *NatsService) Quit() {
	internal.GLog.Info("Quit NatsService ....")
	self.exit <- true
}
