package gateway

import (
	"gateway/src/pb"
	"time"
)

/**
跟前端的心跳处理
*/

func (agent *Agent) HeartBeatHandler(head *pb.ClientCommonHead, request *pb.HeartbeatRequest) {
	agent.Log.Info("HeartBeatHandler %+v", request)
	response := &pb.HeartbeatResponse{Timestamp: time.Now().Unix()}
	// agent.GamePing()
	agent.ReplyClient(head, response)
}
