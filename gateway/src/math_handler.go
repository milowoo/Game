package gateway

import (
	"gateway/src/constants"
	"gateway/src/pb"
)

func (agent *Agent) MatchRequest() {
	if agent.IsMatching {
		agent.Log.Warn("MatchRequest uid %+v is matching", agent.Uid)
		agent.MatchResponse(constants.PLAYER_IS_MATCHING, "player is matching")
		return
	}

	gameInfo := agent.DynamicConfig.GetGameInfo(agent.GameId)
	if gameInfo == nil {
		agent.Log.Warn("MatchRequest uid %+v get game info err gameId %+v", agent.Uid, agent.GameId)
		agent.MatchResponse(constants.INVALID_GAME_ID, "system err")
		return
	}

	var opt = "room"

	if gameInfo.Type == constants.GAME_TYPE_SINGLE {
		if len(agent.RoomId) > 1 {
			agent.Log.Warn("MatchRequest uid %+v get game info err gameId %+v", agent.Uid, agent.GameId)
			agent.MatchResponse(constants.PLAYER_IN_ROOM, "have in room")
			return
		}
	} else if gameInfo.Type == constants.GAME_TYPE_HALL_SINGLE {
		if agent.InHall == 2 {
			agent.Log.Warn("MatchRequest uid %+v get game info err gameId %+v", agent.Uid, agent.GameId)
			agent.MatchResponse(constants.PLAYER_IN_ROOM, "have in room")
			return
		}
		if agent.InHall == 0 {
			opt = "hall"
		}
	} else if gameInfo.Type == constants.GAME_TYPE_HALL_1V1 {
		if agent.InHall == 2 {
			agent.Log.Warn("MatchRequest uid %+v get game info err gameId %+v", agent.Uid, agent.GameId)
			agent.MatchResponse(constants.PLAYER_IN_ROOM, "have in room")
			return
		}
		if agent.InHall == 0 {
			opt = "hall"
		}
	} else if gameInfo.Type == constants.GAME_TYPE_1V1 {
		if len(agent.RoomId) > 1 {
			agent.Log.Warn("MatchRequest uid %+v get game info err gameId %+v", agent.Uid, agent.GameId)
			agent.MatchResponse(constants.PLAYER_IN_ROOM, "have in room")
			return
		}
	}

	//发起匹配请求
	agent.Server.MatchRequest(agent.GameId, agent.Uid, 1, opt)
	agent.IsMatching = true
	agent.MatchResponse(constants.CODE_SUCCESS, "success")
}

func (agent *Agent) MatchResponse(code int32, msg string) {
	binary, err := GetBinary(&pb.ClientMatchResponse{
		Code: code,
		Msg:  msg},
		agent.Log, agent.Config.AgentConfig)
	if err != nil {
		return
	}
	agent.SendBinaryNow(binary)
}

func (agent *Agent) MatchOverResponse(res *pb.MatchOverRes) {
	if agent.GameId != res.GetGameId() {
		agent.Log.Error("MatchResponse uid %+v match game %+v agent game %+v", agent.Uid, res.GetGameId(), agent.GameId)
		return
	}

	gameInfo := agent.DynamicConfig.GetGameInfo(agent.GameId)
	if gameInfo == nil {
		agent.Log.Warn("MatchOverResponse uid %+v get game info err gameId %+v", agent.Uid, agent.GameId)
		agent.MatchResponse(constants.SYSTEM_ERROR, "system err")
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
	agent.ReplyClient(res)
	agent.GameSubject = constants.GetGameSubject(gameInfo.GameId, res.GetGameIp())
}
