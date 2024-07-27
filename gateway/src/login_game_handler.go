package gateway

import (
	"gateway/src/pb"
	"github.com/golang/protobuf/proto"
)

func (agent *Agent) DoLoginReply(code int32, msg string) {
	response := &pb.LoginResponse{
		Code: code,
		Msg:  msg,
		Uid:  agent.Uid}

	protoName := proto.MessageName(response)
	head := &pb.ClientCommonHead{Pid: agent.Pid,
		Sn:        agent.Counter.GetIncrementValue(),
		ProtoName: protoName}

	agent.ReplyClient(head, response)
	agent.Log.Info("uid %v loginReply, %+v, %+v", agent.Uid, code, msg)
}
