package ucenter

import (
	"github.com/golang/protobuf/proto"
	"ucenter/src/constants"
	"ucenter/src/domain"
	"ucenter/src/pb"
)

func (self *HandlerMgr) GetPlayerUID(reply string, request *pb.ApplyUidRequest) {
	//var request pb.ApplyUidRequest
	//err := proto.Unmarshal(msg.Data, &request)
	//if err != nil {
	//	response := &pb.ApplyUidResponse{
	//		Code: constants.INVALID_BODY,
	//		Msg:  "Unmarshal failed",
	//	}
	//	res, _ := proto.Marshal(response)
	//	// 伪同步响应：接收到请求消息后需向响应收件箱发送一条消息作为回应
	//	err = self.NatsPool.Publish(reply, res)
	//	if err != nil {
	//		self.log.Error("SubscribeGetUid reply err %+v", err)
	//	}
	//	return
	//}

	var uid string
	pid := request.GetPid()
	player := self.server.MongoDao.GetPlayer(pid)
	if player == nil {
		uid = self.CreateUid()
		dbPlayer := domain.NewPlayer(pid, uid)
		dbPlayer.Uid = uid
		self.server.MongoDao.InsertPlayer(dbPlayer)
	} else {
		uid = player.Uid
	}

	self.getUidResponse(reply, constants.CODE_SUCCESS, "system err", pid, uid)
}

func (self *HandlerMgr) getUidResponse(reply string, code int32, msg string, pid string, uid string) {
	response := &pb.ApplyUidResponse{
		Code: code,
		Msg:  msg,
		Uid:  uid,
		Pid:  pid,
	}

	resByte, _ := proto.Marshal(response)
	// 伪同步响应：接收到请求消息后需向响应收件箱发送一条消息作为回应
	self.NatsPool.Publish(reply, map[string]interface{}{"res": "ok", "data": string(resByte)})
}
