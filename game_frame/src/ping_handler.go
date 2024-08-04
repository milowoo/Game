package game_frame

import (
	"game_frame/src/constants"
	"game_frame/src/pb"
	"github.com/golang/protobuf/proto"
)

func (self *Room) PingHandler(reply string, head *pb.CommonHead, data []byte) {
	player := self.GetPlayer(head.GetUid())
	if player == nil {
		self.responsePing(reply, head, constants.INVALID_USER_ID, "invalid user id")
		return
	}

	player.OffLineFrameId = self.FrameId
	self.responsePing(reply, head, constants.CODE_SUCCESS, "")
}

func (self *Room) responsePing(reply string, head *pb.CommonHead, code int32, msg string) {
	response := &pb.PingResponse{
		Code: code,
		Msg:  msg,
	}

	head.PbName = proto.MessageName(response)
	self.ResponseGateway(reply, head, response)
}
