package gateway

import (
	"gateway/src/constants"
	"gateway/src/pb"
	"github.com/golang/protobuf/proto"
	"time"
)

func (agent *Agent) LoginHallHandler(head *pb.ClientCommonHead, request *pb.ClientLoginHallRequest) {
	agent.Log.Info("LoginHallHandler uid %+v begin ", agent.Uid)
	if agent.IsMatching {
		agent.Log.Warn("LoginHallHandler uid %+v is matching", agent.Uid)
		agent.LoginHallResponse(head, constants.PLAYER_IS_MATCHING, "player is matching")
		return
	}

	gameInfo := agent.DynamicConfig.GetGameInfo(agent.GameId)
	if gameInfo == nil {
		agent.Log.Warn("LoginHallHandler uid %+v get game info err gameId %+v", agent.Uid, agent.GameId)
		agent.LoginHallResponse(head, constants.SYSTEM_ERROR, "system err")
		return
	}

	if gameInfo.Type != constants.GAME_TYPE_HALL_1V1 && gameInfo.Type != constants.GAME_TYPE_HALL_SINGLE {
		agent.Log.Warn("LoginHallHandler uid %+v invalid request %+v", agent.Uid, agent.GameId)
		agent.LoginHallResponse(head, constants.INVALID_REQUEST, "invalid request")
		return
	}

	if gameInfo.Type == constants.GAME_TYPE_HALL_SINGLE {
		if agent.InHall == 2 {
			agent.Log.Warn("MatchRequest uid %+v get game info err gameId %+v", agent.Uid, agent.GameId)
			agent.LoginHallResponse(head, constants.PLAYER_IN_ROOM, "have in room")
			return
		}
	} else if gameInfo.Type == constants.GAME_TYPE_HALL_SINGLE {
		if agent.InHall == 2 {
			agent.Log.Warn("LoginHallHandler uid %+v get game info err gameId %+v", agent.Uid, agent.GameId)
			agent.LoginHallResponse(head, constants.PLAYER_IN_ROOM, "have in room")
			return
		}
	}

	if agent.LoginHall2Match(head) != nil {
		agent.Log.Warn("LoginHallHandler uid %+v get game info err gameId %+v", agent.Uid, agent.GameId)
		agent.LoginHallResponse(head, constants.SYSTEM_ERROR, "system err")
		return
	}

	agent.InHall = 1

}

func (agent *Agent) LoginHall2Match(head *pb.ClientCommonHead) error {
	request := pb.CreateHallRequest{
		Uid:    agent.Uid,
		GameId: agent.GameId,
	}

	bytes, _ := proto.Marshal(&request)
	var response interface{}

	//发起进入大厅请求
	agent.Log.Info("LonginHall2Match begin %+v", agent.Uid)
	agent.Server.NatsPool.Request(constants.LOGIN_HALL_SUBJECT, string(bytes), &response, 3*time.Second)
	agent.Log.Info("LonginHall2Match response %+v", response)

	dataMap := response.(map[string]interface{})
	resBytes := []byte(dataMap["data"].(string))
	var res pb.CreateHallResponse
	err := proto.Unmarshal(resBytes, &res)
	if err != nil {
		agent.Log.Error("LonginHallRequestMatch err %+v", err)
		return err
	}

	agent.RoomId = res.GetRoomId()
	agent.GameSubject = constants.GetGameSubject(agent.GameId, res.GetGameIp())

	agent.LoginGameHall(head)
	return nil
}

func (agent *Agent) LoginGameHall(head *pb.ClientCommonHead) *pb.LoginHallResponse {
	request := &pb.LoginHallRequest{}

	//发起进入大厅请求
	agent.Log.Info("LoginGameHall begin %+v GameSubject %+v", agent.Uid, agent.GameSubject)
	protoName := proto.MessageName(request)
	response, err := agent.RequestToGame(protoName, request)
	if err != nil {
		agent.Log.Error("LoginGameHall err %+v", err)
		agent.DoLoginReply(constants.SYSTEM_ERROR, "system err")
		return nil
	}

	agent.Log.Info("LoginGameHall response %+v", response)

	var res pb.LoginHallResponse
	proto.Unmarshal(response, &res)

	agent.ReplyClient(head, &pb.ClientLoginHallResponse{Code: constants.CODE_SUCCESS,
		Msg:      "success",
		JumpGame: res.GetJumpGame(),
		RoomId:   res.GetRoomId(),
		Opt:      res.GetOpt(),
	})

	agent.InHall = 1

	return &res
}

func (agent *Agent) LoginHallResponse(head *pb.ClientCommonHead, code int32, msg string) {
	agent.ReplyClient(head, &pb.ClientLoginHallResponse{Code: code, Msg: msg})
}
