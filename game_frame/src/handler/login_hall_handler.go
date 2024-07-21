package handler

import (
	"game_frame/src/constants"
	"game_frame/src/pb"
)

func (self *Room) LoginHall(reply string, head *pb.CommonHead, request *pb.LoginHallRequest) {
	self.isHall = true

	//查询用户的游戏信息存放在内存中
	uid := head.Uid
	if self.IsFirstLogin(uid) {
		player := NewPlayer(head.GetRoomId(), head.GetUid(), head.GetPid(), head.GetHostIp(), false)
		self.Players = append(self.Players, player)
	}

	//记录用户的在大厅信息
	response := &pb.LoginHallResponse{
		Code: constants.CODE_SUCCESS,
		Msg:  "",
	}
	
	self.ResponseGateway(reply, head, response)
}
