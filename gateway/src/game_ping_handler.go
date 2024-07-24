package gateway

import (
	"gateway/src/pb"
	"reflect"
	"time"
)

/**
跟游戏服务的心跳处理
*/

func (agent *Agent) GamePing() {
	request := &pb.PingRequest{
		Timestamp: time.Now().Unix(),
	}

	typ := reflect.TypeOf(request)
	protoName := typ.Elem().Name()
	agent.PublicToGame(protoName, request)
}
