package gateway

import (
	"gateway/src/pb"
	"reflect"
)

func UserExitHandler(agent *Agent, reason int32) {
	if len(agent.GameSubject) < 1 {
		agent.Log.Error("UserExitHandler invalid request ")
		return
	}

	request := &pb.UserExitRequest{
		Reason: reason,
	}

	typ := reflect.TypeOf(request)
	protoName := typ.Elem().Name()
	agent.PublicToGame(protoName, request)
}
