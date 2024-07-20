package handler

import (
	gateway "gateway/src"
	"gateway/src/constants"
	"gateway/src/pb"
)

func CancelMatchRequest(agent *gateway.Agent) {
	if !agent.IsMatching {
		agent.Log.Warn("CancelMatchRequest uid %+v is not in matching", agent.Uid)
		MatchResponse(agent, constants.PLAYER_NO_MATCHING, "player is no in matching")
		return
	}

	//调用match 服务
	agent.Server.CancelMatchRequest(agent.GameId, agent.Uid)

	agent.IsMatching = false
	binary, err := gateway.GetBinary(&pb.ClientCancelMatchResponse{
		Code: 0,
		Msg:  "success"},
		agent.Log, agent.Config.AgentConfig)
	if err != nil {
		return
	}
	agent.SendBinaryNow(binary)

}
