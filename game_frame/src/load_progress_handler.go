package game_frame

import (
	"game_frame/src/constants"
	"game_frame/src/internal"
	"game_frame/src/pb"
	"github.com/golang/protobuf/proto"
)

func (self *Room) LoadProgressHandler(reply string, head *pb.CommonHead, data []byte) {
	var request pb.LoadProgressRequest
	proto.Unmarshal(data, &request)

	internal.GLog.Info("LoadProgressHandler reply: %s request [%+v]", reply, &request)

	//转发进度给对手
	rival := self.GetRival(head.GetUid())
	if rival != nil {
		self.Send2PlayerMessage(rival, &request)
	}

	progressList := make([]*pb.Progress, 0)
	progress := &pb.Progress{
		Uid:      head.Uid,
		Progress: request.GetProgress(),
	}

	progressList = append(progressList, progress)

	rivalProgress := self.getRivalProgress(head.GetUid())
	if rivalProgress != nil {
		progressList = append(progressList, rivalProgress)
	}

	response := &pb.LoadProgressResponse{
		Code: constants.CODE_SUCCESS,
		Msg:  "",
		Data: progressList,
	}

	head.PbName = proto.MessageName(response)
	self.ResponseGateway(reply, head, response)
}

func (self *Room) getRivalProgress(uid string) *pb.Progress {
	player := self.GetRival(uid)
	if player == nil {
		return nil
	}

	return &pb.Progress{
		Uid:      player.Uid,
		Progress: player.LoadingProgress,
	}
}
