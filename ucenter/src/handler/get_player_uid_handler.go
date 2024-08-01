package handler

import (
	"github.com/golang/protobuf/proto"
	"strconv"
	"ucenter/src/constants"
	"ucenter/src/domain"
	"ucenter/src/internal"
	"ucenter/src/pb"
)

func GetPlayerUID(reply string, request *pb.ApplyUidRequest) {
	var uid string
	pid := request.GetPid()
	player := internal.Mongo.GetPlayer(pid)
	if player == nil {
		uid = createUid()
		dbPlayer := domain.NewPlayer(pid, uid)
		dbPlayer.Uid = uid
		internal.Mongo.InsertPlayer(dbPlayer)
	} else {
		uid = player.Uid
	}

	getUidResponse(reply, constants.CODE_SUCCESS, "", pid, uid)
}

func getUidResponse(reply string, code int32, msg string, pid string, uid string) {
	response := &pb.ApplyUidResponse{
		Code: code,
		Msg:  msg,
		Uid:  uid,
		Pid:  pid,
	}

	resByte, _ := proto.Marshal(response)
	// 伪同步响应：接收到请求消息后需向响应收件箱发送一条消息作为回应
	internal.NatsPool.Publish(reply, map[string]interface{}{"res": "ok", "data": string(resByte)})
}

func createUid() string {
	uid, _ := internal.RedisDao.IncrBy("create_uid", 1)
	return strconv.FormatInt(uid+1000, 10)
}
