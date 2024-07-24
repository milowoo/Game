package gateway

import (
	"gateway/src/pb"
	"time"
)

/**
跟前端的心跳处理
*/

func (agent *Agent) HeartBeatReply() {
	binary, err := GetBinary(&pb.HeartbeatResponse{
		Timestamp: time.Now().Unix()},
		agent.Log, agent.Config.AgentConfig)
	if err != nil {
		return
	}
	agent.GamePing()
	agent.SendBinaryNow(binary)
}
