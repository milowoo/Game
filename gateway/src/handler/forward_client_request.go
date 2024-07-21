package handler

import (
	gateway "gateway/src"
	"gateway/src/constants"
	"gateway/src/pb"
	"github.com/gogo/protobuf/proto"
	"github.com/nats-io/go-nats"
	"time"
)

func ForwardClientRequest(agent *gateway.Agent, protoName string, request proto.Message) {
	if len(agent.GameSubject) < 1 {
		agent.Log.Error("ForwardClientRequest uid %+v protoName %+v invalid", agent.Uid, protoName)
		return
	}

	bytes, _ := proto.Marshal(request)

	head := &pb.CommonHead{
		GameId:    agent.GameId,
		Uid:       agent.Uid,
		RoomId:    agent.RoomId,
		Sn:        agent.Counter.GetIncrementValue(),
		ProtoName: protoName,
		HostIp:    gateway.GetHostIp(),
	}

	commonRequest := &pb.GameCommonRequest{
		Head: head,
		Data: bytes,
	}

	commonBytes, _ := proto.Marshal(commonRequest)

	var response interface{}
	err := agent.Server.NatsPool.Request(agent.GameSubject, commonBytes, &response, 3*time.Second)
	if err != nil {
		if agent.RequestGameErrFrame == 0 {
			agent.RequestGameErrFrame = agent.FrameID
		}
		agent.Log.Error("ForwardNeedResponse err %+v", err)
		return
	}
	agent.RequestGameErrFrame = 0
	data, _ := response.(*nats.Msg)
	var res pb.GameCommonResponse
	err = proto.Unmarshal(data.Data, &res)
	if err != nil {
		agent.Log.Error("ForwardNeedResponse err %+v", err)
		return
	}

	if res.Code != constants.CODE_SUCCESS {
		agent.Log.Error("ForwardClientRequest uid %+v protoName %+v err", agent.Uid, protoName)
	} else {
		var response proto.Message
		proto.Unmarshal(data.Data, response)
		agent.ReplyClient(response)
	}

}
