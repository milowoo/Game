package handler

import (
	gateway "gateway/src"
	"gateway/src/pb"
	"reflect"
)

func UserExitHandler(agent *gateway.Agent, reason int32) {
	if len(agent.GameSubject) < 1 {
		agent.Log.Error("UserExitHandler invalid request ")
		return
	}

	request := &pb.UserExitRequest{
		Reason: reason,
	}

	typ := reflect.TypeOf(request)
	protoName := typ.Elem().Name()
	PublicToGame(agent, protoName, request)
}
