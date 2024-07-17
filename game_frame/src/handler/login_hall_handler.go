package handler

import pb "game_frame/src/pb"

func (self *Room) LoginHall(commonRes *pb.GameCommonResponse, reply string, uid string, request *pb.LoginHallRequest) {
	self.isHall = true

	//查询用户的游戏信息存放在内存中

}
