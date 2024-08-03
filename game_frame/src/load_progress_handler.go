package game_frame

import (
	"game_frame/src/internal"
	"game_frame/src/pb"
)

func (self *Room) LoadProgressHandler(reply string, head *pb.CommonHead, request *pb.LoadProgressRequest) {
	internal.GLog.Info("LoadProgressReq reply: %s request %+v", reply, &request)
	//转发进度给对手
	rival := self.GetRival(head.GetUid())
	if rival != nil {
		self.Send2PlayerMessage(rival, request)
	}
}
