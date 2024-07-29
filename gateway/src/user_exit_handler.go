package gateway

import (
	"gateway/src/pb"
	"github.com/golang/protobuf/proto"
)

func UserExitHandler(agent *Agent, reason int32) {
	if len(agent.GameSubject) < 1 {
		agent.Log.Error("UserExitHandler invalid request ")
		return
	}

	request := &pb.UserExitRequest{
		Reason: reason,
	}

	protoName := proto.MessageName(request)
	agent.RequestToGame(protoName, request)
}
