package handler

import (
	gateway "gateway/src"
	"gateway/src/pb"
	"github.com/gogo/protobuf/proto"
	"time"
)

/**
跟游戏服务的心跳处理
*/

func GamePing(agent *gateway.Agent) {
	if len(agent.GameSubject) < 1 {
		return
	}

	request := &pb.PingRequest{
		Timestamp: time.Now().Unix(),
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
		agent.Log.Error("uid %v ping game err %+v", agent.Uid, err)
		return
	}
}
