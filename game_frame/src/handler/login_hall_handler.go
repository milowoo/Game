package handler

import (
	"encoding/json"
	"game_frame/src"
	"game_frame/src/constants"
	"game_frame/src/domain"
	"game_frame/src/pb"
	"time"
)

func (self *Room) LoginHall(reply string, head *pb.CommonHead, request *pb.LoginHallRequest) {
	self.isHall = true

	//查询用户的游戏信息存放在内存中

	//记录用户的在大厅信息
	hostIp := game_frame.GetHostIp()
	hallData := domain.HallData{
		RoomId: head.RoomId,
		Uid:    head.Uid,
		GameId: head.GameId,
		UTime:  time.Now().Unix(),
	}

	redisKey := "game." + hostIp
	data, _ := json.Marshal(hallData)
	self.RedisDao.SAdd(redisKey, data)
	response := &pb.LoginHallResponse{
		Code:   constants.CODE_SUCCESS,
		RoomId: head.RoomId,
	}
	self.ResponseGateway(reply, head, response)
}
