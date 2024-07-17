package handler

import (
	"github.com/golang/protobuf/proto"
	"github.com/nats-io/go-nats"
	"ucenter/src/domain"
	"ucenter/src/pb"
)

func (self *HandlerMgr) GetPlayerUID(reply string, msg *nats.Msg) {
	var request pb.ApplyUidRequest
	err := proto.Unmarshal(msg.Data, &request)
	if err != nil {
		response := &pb.ApplyUidResponse{
			Code: 101,
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
	player, err := self.server.MongoDao.GetPlayerByPid(pid)
	if err != nil {
		self.log.Error("fail to get player by sign(%+v): %+v", pid, err)
	}
	var uid string
	if player == nil {
		uid = self.GreateUid()
		dbPlayer := domain.NewPlayer(uid)
		dbPlayer.Pid = pid
		self.server.MongoDao.InsertPlayer(dbPlayer)
	} else {
		uid = player.Uid.String()
	}

	response := &pb.ApplyUidResponse{
		Code: 0,
		Msg:  "success",
		Uid:  uid,
		Pid:  pid,
	}

	res, _ := proto.Marshal(response)
	// 伪同步响应：接收到请求消息后需向响应收件箱发送一条消息作为回应
	err = self.NatsPool.Publish(reply, res)
	if err != nil {
		self.log.Error("SubscribeGetUid reply err %+v", err)
	}
}
