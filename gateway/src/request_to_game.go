package gateway

import (
	"gateway/src/pb"
	"github.com/golang/protobuf/proto"
	"time"
)

func (agent *Agent) RequestToGame(protoName string, protoMsg proto.Message) ([]byte, error) {
	agent.Log.Info("PublicToGame protoName %+v", protoName)
	if len(agent.GameSubject) < 1 {
		agent.Log.Error("PublicToGame invalid request %+v", protoName)
		return nil, nil
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
	err := agent.Server.NatsPool.Request(agent.GameSubject, string(comBytes), &response, 3*time.Second)
	if err != nil {
		agent.Log.Error("uid %v protoName %v public to game game err %+v", agent.Uid, protoName, err)
		if agent.RequestGameErrFrame == 0 {
			agent.RequestGameErrFrame = agent.FrameID
		}

		return nil, err
	}
	agent.RequestGameErrFrame = 0

	dataMap := response.(map[string]interface{})
	resBytes := []byte(dataMap["data"].(string))
	return resBytes, nil
}
