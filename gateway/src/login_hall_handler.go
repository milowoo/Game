package gateway

import (
	"gateway/src/constants"
	"gateway/src/pb"
	"github.com/gogo/protobuf/proto"
	"github.com/nats-io/go-nats"
	"time"
)

func (agent *Agent) LoginHallRequest() {
	if agent.IsMatching {
		agent.Log.Warn("LoginHallRequest uid %+v is matching", agent.Uid)
		agent.LoginHallResponse(constants.PLAYER_IS_MATCHING, "player is matching")
		return
	}

	gameInfo := agent.DynamicConfig.GetGameInfo(agent.GameId)
	if gameInfo == nil {
		agent.Log.Warn("LoginHallRequest uid %+v get game info err gameId %+v", agent.Uid, agent.GameId)
		agent.LoginHallResponse(constants.SYSTEM_ERROR, "system err")
		return
	}

	if gameInfo.Type != constants.GAME_TYPE_HALL_1V1 && gameInfo.Type != constants.GAME_TYPE_HALL_SINGLE {
		agent.Log.Warn("LoginHallRequest uid %+v invalid request %+v", agent.Uid, agent.GameId)
		agent.LoginHallResponse(constants.INVALID_REQUEST, "invalid request")
		return
	}

	if gameInfo.Type == constants.GAME_TYPE_HALL_SINGLE {
		if agent.InHall == 2 {
			agent.Log.Warn("MatchRequest uid %+v get game info err gameId %+v", agent.Uid, agent.GameId)
			agent.LoginHallResponse(constants.PLAYER_IN_ROOM, "have in room")
			return
		}
	} else if gameInfo.Type == constants.GAME_TYPE_HALL_SINGLE {
		if agent.InHall == 2 {
			agent.Log.Warn("MatchRequest uid %+v get game info err gameId %+v", agent.Uid, agent.GameId)
			agent.LoginHallResponse(constants.PLAYER_IN_ROOM, "have in room")
			return
		}
	}

	if agent.LonginHall2Match() != nil {
		agent.Log.Warn("LoginHallRequest uid %+v get game info err gameId %+v", agent.Uid, agent.GameId)
		agent.LoginHallResponse(constants.SYSTEM_ERROR, "system err")
		return
	}

	agent.LoginHallResponse(constants.CODE_SUCCESS, "success")
}

func (agent *Agent) LonginHall2Match() error {
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

func (agent *Agent) LoginHall2Game() {
	//todo 111
}

func (agent *Agent) LoginHallResponse(code int32, msg string) {
	binary, err := GetBinary(&pb.ClientLoginHallResponse{
		Code: code,
		Msg:  msg},
		agent.Log, agent.Config.AgentConfig)
	if err != nil {
		return
	}
	agent.SendBinaryNow(binary)
}
