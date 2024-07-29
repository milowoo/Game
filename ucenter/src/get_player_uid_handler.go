package ucenter

import (
	"github.com/golang/protobuf/proto"
	"ucenter/src/constants"
	"ucenter/src/domain"
	"ucenter/src/pb"
)

func (self *HandlerMgr) GetPlayerUID(reply string, request *pb.ApplyUidRequest) {
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

	self.getUidResponse(reply, constants.CODE_SUCCESS, "", pid, uid)
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
