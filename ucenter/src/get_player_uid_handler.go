package ucenter

import (
	"github.com/golang/protobuf/proto"
	"github.com/nats-io/go-nats"
	"ucenter/src/constants"
	"ucenter/src/domain"
	"ucenter/src/pb"
)

func (self *HandlerMgr) GetPlayerUID(reply string, msg *nats.Msg) {
	var request pb.ApplyUidRequest
	err := proto.Unmarshal(msg.Data, &request)
	if err != nil {
		response := &pb.ApplyUidResponse{
			Code: constants.INVALID_BODY,
			Msg:  "Unmarshal failed",
		}
		res, _ := proto.Marshal(response)
		// 伪同步响应：接收到请求消息后需向响应收件箱发送一条消息作为回应
		err = self.NatsPool.Publish(reply, res)
		if err != nil {
			self.log.Error("SubscribeGetUid reply err %+v", err)
		}
		return
	}
	pid := request.GetPid()
	player, err := self.server.MongoDao.GetPlayer(pid)
	if err != nil {
		self.log.Error("fail to get player by sign(%+v): %+v", pid, err)
	}
	var uid string
	if player == nil {
		uid = self.GreateUid()
		dbPlayer := domain.NewPlayer(pid, uid)
		dbPlayer.Uid = pid
		self.server.MongoDao.InsertPlayer(dbPlayer)
	} else {
		uid = player.Uid
	}

	response := &pb.ApplyUidResponse{
		Code: constants.CODE_SUCCESS,
		Msg:  "success",
		Uid:  uid,
		Pid:  pid,
	}

	res, _ := proto.Marshal(response)
	// 伪同步响应：接收到请求消息后需向响应收件箱发送一条消息作为回应
	err = self.NatsPool.Publish(reply, res)
	if err != nil {
		self.log.Error("SubscribeGetUid reply pid %+v err %+v ", request.Pid, err)
	}
}
