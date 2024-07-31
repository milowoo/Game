package game_frame

import (
	"game_frame/src/constants"
	"game_frame/src/pb"
	"github.com/golang/protobuf/proto"
)

func (self *Room) LoginHall(reply string, head *pb.CommonHead, data []byte) {
	self.isHall = true
	self.Log.Info("LoginHall reply: %s", reply)
	var request pb.LoginHallRequest
	err := proto.Unmarshal(data, &request)
	if err != nil {
		self.Log.Info("LoginHall proto.Unmarshal err: %s", err.Error())
	}

	//查询用户的游戏信息存放在内存中
	uid := head.Uid
	if self.IsFirstLogin(uid) {
		player := NewPlayer(head.GetRoomId(), head.GetUid(), head.GetPid(), head.GetGatewayIp(), false)
		self.Players = append(self.Players, player)
	}

	//记录用户的在大厅信息
	response := &pb.LoginHallResponse{
		Code: constants.CODE_SUCCESS,
		Msg:  "",
	}

	self.Log.Info("LoginHall reply: uid %+v success", uid)
	head.PbName = proto.MessageName(response)
	self.ResponseGateway(reply, head, response)
}
