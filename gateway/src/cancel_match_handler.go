package gateway

import (
	"gateway/src/constants"
	"gateway/src/internal"
	"gateway/src/pb"
)

func (agent *Agent) CancelMatchHandler(head *pb.ClientCommonHead, request *pb.CancelMatchRequest) {
	if !agent.IsMatching {
		internal.GLog.Warn("CancelMatchRequest uid %+v is not in matching", agent.Uid)
		agent.MatchResponse(head, constants.PLAYER_NO_MATCHING, "player is no in matching")
		return
	}

	//调用match 服务
	agent.Server.CancelMatchRequest(agent.GameId, agent.Uid)

	agent.IsMatching = false
	agent.CancelMatchResponse(head, constants.CODE_SUCCESS, "")

}

func (agent *Agent) CancelMatchResponse(head *pb.ClientCommonHead, code int32, msg string) {
	response := &pb.ClientCancelMatchResponse{
		Code: code,
		Msg:  msg}

	agent.ReplyClient(head, response)
}
