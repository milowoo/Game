package gateway

import (
	"gateway/src/internal"
	"gateway/src/pb"
	"time"
)

/**
跟前端的心跳处理
*/

func (agent *Agent) HeartBeatHandler(head *pb.ClientCommonHead, request *pb.HeartbeatRequest) {
	internal.GLog.Info("HeartBeatHandler %+v", request)
	response := &pb.HeartbeatResponse{Timestamp: time.Now().UnixMilli()}
	// agent.GamePing()
	agent.ReplyClient(head, response)
}
