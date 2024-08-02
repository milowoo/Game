package gateway

import (
	"gateway/src/constants"
	"gateway/src/internal"
	"gateway/src/pb"
	"github.com/golang/protobuf/proto"
)

func (agent *Agent) MatchHandler(head *pb.ClientCommonHead, request *pb.ClientMatchRequest) {
	if agent.IsMatching {
		internal.GLog.Warn("MatchRequest uid %+v is matching", agent.Uid)
		agent.MatchResponse(head, constants.PLAYER_IS_MATCHING, "player is matching")
		return
	}

	gameInfo := agent.DynamicConfig.GetGameInfo(agent.GameId)
	if gameInfo == nil {
		internal.GLog.Warn("MatchRequest uid %+v get game info err gameId %+v", agent.Uid, agent.GameId)
		agent.MatchResponse(head, constants.INVALID_GAME_ID, "system err")
		return
	}

	var opt = "room"

	if gameInfo.Type == constants.GAME_TYPE_SINGLE {
		if len(agent.RoomId) > 1 {
			internal.GLog.Warn("MatchRequest uid %+v get game info err gameId %+v", agent.Uid, agent.GameId)
			agent.MatchResponse(head, constants.PLAYER_IN_ROOM, "have in room")
			return
		}
	} else if gameInfo.Type == constants.GAME_TYPE_HALL_SINGLE {
		if agent.InHall == 2 {
			internal.GLog.Warn("MatchRequest uid %+v get game info err gameId %+v", agent.Uid, agent.GameId)
			agent.MatchResponse(head, constants.PLAYER_IN_ROOM, "have in room")
			return
		}
		if agent.InHall == 0 {
			opt = "hall"
		}
	} else if gameInfo.Type == constants.GAME_TYPE_HALL_1V1 {
		if agent.InHall == 2 {
			internal.GLog.Warn("MatchRequest uid %+v get game info err gameId %+v", agent.Uid, agent.GameId)
			agent.MatchResponse(head, constants.PLAYER_IN_ROOM, "have in room")
			return
		}
		if agent.InHall == 0 {
			opt = "hall"
		}
	} else if gameInfo.Type == constants.GAME_TYPE_1V1 {
		if len(agent.RoomId) > 1 {
			internal.GLog.Warn("MatchRequest uid %+v get game info err gameId %+v", agent.Uid, agent.GameId)
			agent.MatchResponse(head, constants.PLAYER_IN_ROOM, "have in room")
			return
		}
	}

	//发起匹配请求
	agent.Server.MatchRequest(agent.GameId, agent.Uid, 1, opt)
	agent.IsMatching = true
	agent.MatchResponse(head, constants.CODE_SUCCESS, "success")
}

func (agent *Agent) MatchResponse(head *pb.ClientCommonHead, code int32, msg string) {
	response := &pb.ClientMatchResponse{
		Code: code,
		Msg:  msg}

	agent.ReplyClient(head, response)
}

func (agent *Agent) MatchOverResponse(res *pb.MatchOverRes) {
	if agent.GameId != res.GetGameId() {
		internal.GLog.Error("MatchResponse uid %+v match game %+v agent game %+v", agent.Uid, res.GetGameId(), agent.GameId)
		return
	}

	protoName := proto.MessageName(res)

	head := &pb.ClientCommonHead{Pid: agent.Pid,
		Sn:        agent.Counter.GetIncrementValue(),
		ProtoName: protoName}

	gameInfo := agent.DynamicConfig.GetGameInfo(agent.GameId)
	if gameInfo == nil {
		internal.GLog.Warn("MatchOverResponse uid %+v get game info err gameId %+v", agent.Uid, agent.GameId)
		agent.MatchResponse(head, constants.SYSTEM_ERROR, "system err")
		return
	}

	if gameInfo.Type == constants.GAME_TYPE_SINGLE || gameInfo.Type == constants.GAME_TYPE_1V1 {
		agent.RoomId = res.GetRoomId()
	} else if gameInfo.Type == constants.GAME_TYPE_HALL_SINGLE || gameInfo.Type == constants.GAME_TYPE_HALL_1V1 {
		if agent.InHall == 0 {
			agent.InHall = 1
		} else {
			agent.InHall = 2
		}

		agent.RoomId = res.GetRoomId()
	}

	agent.IsMatching = false

	//通知客户端
	agent.ReplyClient(head, res)
	agent.GameSubject = constants.GetGameSubject(gameInfo.GameId, res.GetGameIp())
}
