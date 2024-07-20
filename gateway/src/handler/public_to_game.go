package handler

import (
	gateway "gateway/src"
	"gateway/src/pb"
	"github.com/gogo/protobuf/proto"
	"time"
)

func PublicToGame(agent *gateway.Agent, protoName string, protoMsg proto.Message) {
	if len(agent.GameSubject) < 1 {
		agent.Log.Error("PublicToGame invalid request %+v", protoName)
		return
	}

	head := &pb.CommonHead{
		GameId:    agent.GameId,
		Uid:       agent.Uid,
		RoomId:    agent.RoomId,
		Sn:        agent.Counter.GetIncrementValue(),
		Timestamp: time.Now().Unix(),
		ProtoName: protoName,
	}

	bytes, _ := proto.Marshal(protoMsg)
	commonRequest := &pb.GameCommonRequest{
		Head: head,
		Data: bytes,
	}

	comBytes, _ := proto.Marshal(commonRequest)
	var response interface{}
	err := agent.Server.NatsPool.Request(agent.GameSubject, comBytes, &response, 1*time.Second)
	if err != nil {
		agent.Log.Error("uid %v protoName %v public to game game err %+v", agent.Uid, protoName, err)
		if agent.RequestGameErrFrame == 0 {
			agent.RequestGameErrFrame = agent.FrameID
		}

		return
	}
	agent.RequestGameErrFrame = 0
}
