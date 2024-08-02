package gateway

import (
	"gateway/src/constants"
	"gateway/src/internal"
	"gateway/src/pb"
	"github.com/golang/protobuf/proto"
	"time"
)

func (agent *Agent) ForwardClientRequest(client *pb.ClientCommonHead, request proto.Message) {
	protoName := client.GetProtoName()
	if len(agent.GameSubject) < 1 {
		internal.GLog.Error("ForwardClientRequest uid %+v protoName %+v invalid", agent.Uid, protoName)
		return
	}

	bytes, _ := proto.Marshal(request)

	head := &pb.CommonHead{
		GameId:    agent.GameId,
		Uid:       agent.Uid,
		RoomId:    agent.RoomId,
		Sn:        agent.Counter.GetIncrementValue(),
		PbName:    protoName,
		GatewayIp: GetHostIp(),
	}

	commonRequest := &pb.GameCommonRequest{
		Head: head,
		Data: string(bytes),
	}

	commonBytes, _ := proto.Marshal(commonRequest)

	var response interface{}
	err := internal.NatsPool.Request(agent.GameSubject, string(commonBytes), &response, 3*time.Second)
	if err != nil {
		if agent.RequestGameErrFrame == 0 {
			agent.RequestGameErrFrame = agent.FrameID
		}
		internal.GLog.Error("ForwardNeedResponse err %+v", err)
		return
	}
	agent.RequestGameErrFrame = 0

	dataMap := response.(map[string]interface{})
	commonResBytes := []byte(dataMap["data"].(string))
	var res pb.GameCommonResponse
	err = proto.Unmarshal(commonResBytes, &res)
	if err != nil {
		internal.GLog.Error("ForwardNeedResponse err %+v", err)
		return
	}

	if res.Code != constants.CODE_SUCCESS {
		internal.GLog.Error("ForwardClientRequest uid %+v protoName %+v err", agent.Uid, protoName)
	} else {
		var response proto.Message
		proto.Unmarshal([]byte(res.Data), response)
		agent.ReplyClient(client, response)
	}
}
