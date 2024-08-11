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
	progressList := make([]*pb.Progress, 0)

	internal.GLog.Info("LoadProgressHandler reply: %+v request [%+v]", reply, &request)

	player := self.GetPlayer(head.GetUid())
	if player == nil {
		self.responseLoadProgress(reply, head, progressList, constants.INVALID_USER_ID, "invalid user id")
		return
	}

	if player.LoadingProgress >= request.GetProgress() {
		self.responseLoadProgress(reply, head, progressList, constants.INVALID_PROGRESS, "invalid progress")
		return
	}

	if self.State != constants.ROOM_STATE_LOAD {
		self.responseLoadProgress(reply, head, progressList, constants.GAME_IS_RUNING, "game is running")
		return
	}

	//转发进度给对手
	rival := self.GetRival(head.GetUid())
	if rival != nil {
		self.Send2PlayerMessage(rival, &request)
	}

	progress := &pb.Progress{
		Uid:      head.Uid,
		Progress: request.GetProgress(),
	}

	progressList = append(progressList, progress)

	rivalData := int32(0)
	rivalProgress := self.getRivalProgress(head.GetUid())
	if rivalProgress != nil {
		progressList = append(progressList, rivalProgress)
		rivalData = rivalProgress.GetProgress()
	}

	self.responseLoadProgress(reply, head, progressList, constants.CODE_SUCCESS, "")

	// 进度都达到100， 游戏开始
	if rivalData >= 100 && request.GetProgress() >= 100 {
		self.GameRunFrameId = self.FrameId
		self.SetState(constants.ROOM_STATE_START)
	}
}

func (self *Room) responseLoadProgress(reply string, head *pb.CommonHead,
	progressList []*pb.Progress,
	code int32, msg string) {
	response := &pb.LoadProgressResponse{
		Code: code,
		Msg:  msg,
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
