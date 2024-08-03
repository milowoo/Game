package game_frame

import (
	"game_frame/src/constants"
	"game_frame/src/internal"
	"game_frame/src/pb"
	"github.com/golang/protobuf/proto"
)

func (self *Room) LoginHallHandler(reply string, head *pb.CommonHead, request *pb.LoginHallRequest) {
	self.isHall = true
	internal.GLog.Info("LoginHall reply: %s request %+v", reply, &request)
	//var request pb.LoginHallRequest
	//err := proto.Unmarshal(data, &request)
	//if err != nil {
	//	internal.GLog.Info("LoginHall proto.Unmarshal err: %s", err.Error())
	//}

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

	internal.GLog.Info("LoginHall reply: uid %+v success", uid)
	head.PbName = proto.MessageName(response)
	self.ResponseGateway(reply, head, response)
}
