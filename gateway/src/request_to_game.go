package gateway

import (
	"github.com/golang/protobuf/proto"
	"strconv"
	"time"
)

func (agent *Agent) RequestToGame(protoName string, protoMsg proto.Message) ([]byte, error) {
	agent.Log.Info("PublicToGame protoName %+v", protoName)
	if len(agent.GameSubject) < 1 {
		agent.Log.Error("PublicToGame invalid request %+v", protoName)
		return nil, nil
	}

	bytes, _ := proto.Marshal(protoMsg)

	head := map[string]interface{}{
		"gameId":    agent.GameId,
		"uid":       agent.Uid,
		"roomId":    agent.RoomId,
		"pid":       agent.Pid,
		"sn":        strconv.FormatInt(agent.Counter.GetIncrementValue(), 10),
		"timestamp": strconv.FormatInt(time.Now().Unix(), 10),
		"pbName":    protoName,
		"gatewayIp": GetHostIp(),
		"data":      string(bytes),
	}

	var response interface{}
	agent.Log.Info("RequestToGame uid %+v roomId %+v protoName %+v HostIp %+v",
		agent.Uid, agent.RoomId, protoName, GetHostIp())
	err := agent.Server.NatsPool.Request(agent.GameSubject, head, &response, 3*time.Second)
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
