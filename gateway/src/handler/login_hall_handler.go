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
		LoginHallResponse(agent, 102, "system err")
		return
	}

	if gameInfo.Type != constants.GAME_TYPE_HALL_1V1 && gameInfo.Type != constants.GAME_TYPE_HALL_SINGLE {
		agent.Log.Warn("LoginHallRequest uid %+v invalid request %+v", agent.Uid, agent.GameId)
		LoginHallResponse(agent, 103, "invalid request")
		return
	}

	if gameInfo.Type == constants.GAME_TYPE_HALL_SINGLE {
		if agent.InHall == 2 {
			agent.Log.Warn("MatchRequest uid %+v get game info err gameId %+v", agent.Uid, agent.GameId)
			LoginHallResponse(agent, 102, "have in room")
			return
		}
	} else if gameInfo.Type == constants.GAME_TYPE_HALL_SINGLE {
		if agent.InHall == 2 {
			agent.Log.Warn("MatchRequest uid %+v get game info err gameId %+v", agent.Uid, agent.GameId)
			LoginHallResponse(agent, 102, "have in room")
			return
		}
	}

	reqeust := pb.LoginHallRequest{
		Uid:    agent.Uid,
		GameId: agent.GameId,
		Pid:    agent.Pid,
	}

	bytes, _ := proto.Marshal(&reqeust)
	var response interface{}

	//发起进入大厅请求 todo
	agent.Server.NatsPool.Request(agent.GameId, bytes, response, 3*time.Second)
	data, _ := response.(*nats.Msg)
	var res pb.LoginHallResponse
	err := proto.Unmarshal(data.Data, &res)
	if err != nil {
		agent.Log.Error("LoginHallRequest err %+v", err)
		return
	}

	agent.RoomId = res.GetRoomId()
	LoginHallResponse(agent, 0, "success")
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
