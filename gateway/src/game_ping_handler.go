package gateway

import (
	"gateway/src/pb"
	"github.com/golang/protobuf/proto"
	"time"
)

/**
跟游戏服务的心跳处理
*/

func (agent *Agent) GamePing() {
	request := &pb.PingRequest{
		Timestamp: time.Now().UnixMilli(),
	}

	protoName := proto.MessageName(request)
	agent.RequestToGame(protoName, request)
}
