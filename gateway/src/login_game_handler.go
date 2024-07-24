package gateway

import (
	"gateway/src/pb"
	"time"
)

func (agent *Agent) DoLoginReply(code int32, msg string, reconnect int32) {
	binary, err := GetBinary(&pb.LoginResponse{
		Code:      code,
		Msg:       msg,
		Uid:       agent.Uid,
		RoomId:    agent.RoomId,
		Reconnect: reconnect},
		agent.Log, agent.Config.AgentConfig)
	if err != nil {
		return
	}
	agent.SendBinaryNow(binary)
	if code != 200 {
		return
	}

	agent.Log.Debug("uid %v loginReply, %t, %s", agent.Uid, code, msg)
	agent.DelayDisconnect(time.Second * 5)
}
