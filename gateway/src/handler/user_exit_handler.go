package handler

import (
	gateway "gateway/src"
	"gateway/src/pb"
	"github.com/gogo/protobuf/proto"
	"time"
)

func UserExitHandler(agent *gateway.Agent, reason int32) {
	if len(agent.GameSubject) < 1 {
		return
	}

	request := &pb.UserExitRequest{
		Reason: reason,
	}
	bytes, _ := proto.Marshal(request)

	commonRequest := &pb.GameCommonRequest{
		GameId: agent.GameId,
		Uid:    agent.Uid,
		RoomId: agent.RoomId,
		Data:   bytes,
	}

	comBytes, _ := proto.Marshal(commonRequest)
	var response interface{}
	err := agent.Server.NatsPool.Request(agent.GameSubject, comBytes, &response, 1*time.Second)
	if err != nil {
		agent.Log.Error("uid %v exit game err %+v", agent.Uid, err)
		return
	}
}
