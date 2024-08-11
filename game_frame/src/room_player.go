package game_frame

import (
	"game_frame/src/constants"
	"game_frame/src/domain"
	"game_frame/src/internal"
	"game_frame/src/pb"
	"github.com/golang/protobuf/proto"
	"reflect"
	"time"
)

func NewPlayer(roomId string, uid string, pid string, hostIp string, isAi bool) *domain.GamePlayer {
	p := &domain.GamePlayer{
		Uid:       uid,
		Pid:       pid,
		RoomId:    roomId,
		HallId:    "",
		GatewayIp: hostIp,

		IsAi:            isAi,
		IsNewPlayer:     false,
		TotalUseTime:    0,
		OffLineFrameId:  0,
		Score:           0,
		LoadingProgress: 0,
	}

	return p
}

func (self *Room) Send2PlayerMessage(player *domain.GamePlayer, msg proto.Message) {
	if player.IsAi {
		return
	}

	typ := reflect.TypeOf(msg)
	protoName := typ.Elem().Name()
	bytes, _ := proto.Marshal(msg)

	head := &pb.PushHead{
		Uid:       player.Uid,
		GameId:    internal.GameId,
		Pid:       player.Pid,
		Timestamp: time.Now().UnixMilli(),
		RoomId:    player.RoomId,
		PbName:    protoName,
		Sn:        self.Counter.GetIncrementValue(),
	}

	data := &pb.GamePushMessage{
		Head: head,
		Data: string(bytes),
	}
	res, _ := proto.Marshal(data)
	internal.NatsPool.Publish(constants.GetGamePushDataSubject(player.GatewayIp), map[string]interface{}{"res": "ok", "data": string(res)})
}

/*
*
获取竞争对手
*/
func (self *Room) GetRival(uid string) *domain.GamePlayer {
	for _, player := range self.Players {
		if player.Uid != uid {
			return player
		}
	}

	return nil
}

func (self *Room) getRivalUid(uid string) string {
	for _, player := range self.Players {
		if player.Uid != uid {
			return player.Uid
		}
	}

	return ""
}

func (self *Room) GetPlayer(uid string) *domain.GamePlayer {
	for _, player := range self.Players {
		if player.Uid == uid {
			return player
		}
	}

	return nil
}

func (self *Room) IsFirstLogin(uid string) bool {
	if len(self.Players) == 0 {
		return true
	}
	for _, player := range self.Players {
		if player.Uid == uid {
			return false
		}
	}

	return true
}
