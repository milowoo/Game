package handler

import (
	gateway "gateway/src"
	"gateway/src/constants"
	"gateway/src/pb"
	"github.com/gogo/protobuf/proto"
	"github.com/nats-io/go-nats"
	"time"
)

func LoginHallRequest(agent *gateway.Agent) {
	if agent.IsMatching {
		agent.Log.Warn("LoginHallRequest uid %+v is matching", agent.Uid)
		LoginHallResponse(agent, constants.PLAYER_IS_MATCHING, "player is matching")
		return
	}

	gameInfo := agent.DynamicConfig.GetGameInfo(agent.GameId)
	if gameInfo == nil {
		agent.Log.Warn("LoginHallRequest uid %+v get game info err gameId %+v", agent.Uid, agent.GameId)
		LoginHallResponse(agent, constants.SYSTEM_ERROR, "system err")
		return
	}

	if gameInfo.Type != constants.GAME_TYPE_HALL_1V1 && gameInfo.Type != constants.GAME_TYPE_HALL_SINGLE {
		agent.Log.Warn("LoginHallRequest uid %+v invalid request %+v", agent.Uid, agent.GameId)
		LoginHallResponse(agent, constants.INVALID_REQUEST, "invalid request")
		return
	}

	if gameInfo.Type == constants.GAME_TYPE_HALL_SINGLE {
		if agent.InHall == 2 {
			agent.Log.Warn("MatchRequest uid %+v get game info err gameId %+v", agent.Uid, agent.GameId)
			LoginHallResponse(agent, constants.PLAYER_IN_ROOM, "have in room")
			return
		}
	} else if gameInfo.Type == constants.GAME_TYPE_HALL_SINGLE {
		if agent.InHall == 2 {
			agent.Log.Warn("MatchRequest uid %+v get game info err gameId %+v", agent.Uid, agent.GameId)
			LoginHallResponse(agent, constants.PLAYER_IN_ROOM, "have in room")
			return
		}
	}

	if LonginHall2Match(agent) != nil {
		agent.Log.Warn("LoginHallRequest uid %+v get game info err gameId %+v", agent.Uid, agent.GameId)
		LoginHallResponse(agent, constants.SYSTEM_ERROR, "system err")
		return
	}

	LoginHallResponse(agent, 0, "success")
}

func LonginHall2Match(agent *gateway.Agent) error {
	reqeust := pb.CreateHallRequest{
		Uid:    agent.Uid,
		GameId: agent.GameId,
	}

	bytes, _ := proto.Marshal(&reqeust)
	var response interface{}

	//发起进入大厅请求
	agent.Server.NatsPool.Request(agent.GameId, bytes, response, 2*time.Second)
	data, _ := response.(*nats.Msg)
	var res pb.CreateHallResponse
	err := proto.Unmarshal(data.Data, &res)
	if err != nil {
		agent.Log.Error("LonginHallRequestMatch err %+v", err)
		return err
	}

	agent.RoomId = res.GetRoomId()
	agent.GameSubject = constants.GetGameSubject(agent.GameId, res.GetGameIp())
	return nil
}

func LoginHall2Game(agent *gateway.Agent) {

}

func LoginHallResponse(agent *gateway.Agent, code int32, msg string) {
	binary, err := gateway.GetBinary(&pb.ClientLoginHallResponse{
		Code: code,
		Msg:  msg},
		agent.Log, agent.Config.AgentConfig)
	if err != nil {
		return
	}
	agent.SendBinaryNow(binary)
}
