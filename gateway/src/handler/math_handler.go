package handler

import (
	gateway "gateway/src"
	"gateway/src/pb"
)

func MatchRequest(agent *gateway.Agent) {
	if agent.IsMatching {
		agent.Log.Warn("MatchRequest uid %+v is matching", agent.Uid)
		MatchResponse(agent, 101, "player is matching")
		return
	}

	gameInfo := agent.DynamicConfig.GetGameInfo(agent.GameId)
	if gameInfo == nil {
		agent.Log.Warn("MatchRequest uid %+v get game info err gameId %+v", agent.Uid, agent.GameId)
		MatchResponse(agent, 102, "system err")
		return
	}

	if gameInfo.Status != 1 {
		agent.Log.Warn("MatchRequest uid %+v game status %+v err gameId %+v", agent.Uid, gameInfo.Status, agent.GameId)
		MatchResponse(agent, 102, "system err")
		return
	}

	gameType := gateway.GameType(gameInfo.Type)
	var opt = "room"

	if gameType == gateway.Single {
		if len(agent.RoomId) > 1 {
			agent.Log.Warn("MatchRequest uid %+v get game info err gameId %+v", agent.Uid, agent.GameId)
			MatchResponse(agent, 102, "have in room")
			return
		}
	} else if gameType == gateway.HallAndSingle {
		if agent.InHall == 2 {
			agent.Log.Warn("MatchRequest uid %+v get game info err gameId %+v", agent.Uid, agent.GameId)
			MatchResponse(agent, 102, "have in room")
			return
		}
		if agent.InHall == 0 {
			opt = "hall"
		}
	} else if gameType == gateway.HallAnd1V1 {
		if agent.InHall == 2 {
			agent.Log.Warn("MatchRequest uid %+v get game info err gameId %+v", agent.Uid, agent.GameId)
			MatchResponse(agent, 102, "have in room")
			return
		}
		if agent.InHall == 0 {
			opt = "hall"
		}
	} else if gameType == gateway.OneVOne {
		if len(agent.RoomId) > 1 {
			agent.Log.Warn("MatchRequest uid %+v get game info err gameId %+v", agent.Uid, agent.GameId)
			MatchResponse(agent, 102, "have in room")
			return
		}
	}

	//发起匹配请求
	agent.Server.MatchRequest(agent.GameId, agent.Uid, 1, opt)
	agent.IsMatching = true
	MatchResponse(agent, 0, "success")
}

func MatchResponse(agent *gateway.Agent, code int32, msg string) {
	binary, err := gateway.GetBinary(&pb.ClientMatchResponse{
		Code: code,
		Msg:  msg},
		agent.Log, agent.Config.AgentConfig)
	if err != nil {
		return
	}
	agent.SendBinaryNow(binary)
}

func MatchOverResponse(agent *gateway.Agent, res *pb.MatchOverRes) {
	if agent.GameId != res.GetGameId() {
		agent.Log.Error("MatchResponse uid %+v match game %+v agent game %+v", agent.Uid, res.GetGameId(), agent.GameId)
		return
	}

	gameInfo := agent.DynamicConfig.GetGameInfo(agent.GameId)
	if gameInfo == nil {
		agent.Log.Warn("MatchOverResponse uid %+v get game info err gameId %+v", agent.Uid, agent.GameId)
		MatchResponse(agent, 102, "system err")
		return
	}

	gameType := gateway.GameType(gameInfo.Type)
	if gameType == gateway.Single || gameType == gateway.OneVOne {
		agent.RoomId = res.GetRoomId()
	} else if gameType == gateway.HallAndSingle || gameType == gateway.HallAnd1V1 {
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
	agent.GameSubject = "game." + agent.GameId + "." + agent.RoomId
}
