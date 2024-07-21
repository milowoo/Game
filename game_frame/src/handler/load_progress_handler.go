package handler

import "game_frame/src/pb"

func (self *Room) LoadProgressReq(reply string, head *pb.CommonHead, request *pb.LoadProgressRequest) {
	//转发进度给对手
	rival := self.GetRival(head.GetUid())
	if rival != nil {
		self.Send2PlayerMessage(rival, request)
	}
}
