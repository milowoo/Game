package handler

import (
	gateway "gateway/src"
	"gateway/src/pb"
	"reflect"
	"time"
)

/**
跟游戏服务的心跳处理
*/

func GamePing(agent *gateway.Agent) {
	request := &pb.PingRequest{
		Timestamp: time.Now().Unix(),
	}

	typ := reflect.TypeOf(request)
	protoName := typ.Elem().Name()
	PublicToGame(agent, protoName, request)
}
