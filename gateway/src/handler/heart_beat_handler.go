package handler

import (
	gateway "gateway/src"
	"gateway/src/pb"
	"time"
)

/**
跟前端的心跳处理
*/

func HeartBeatReply(agent *gateway.Agent) {
	binary, err := gateway.GetBinary(&pb.HeartbeatResponse{
		Timestamp: time.Now().Unix()},
		agent.Log, agent.Config.AgentConfig)
	if err != nil {
		return
	}
	GamePing(agent)
	agent.SendBinaryNow(binary)
}
