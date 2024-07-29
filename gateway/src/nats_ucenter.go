package gateway

import (
	"gateway/src/constants"
	"gateway/src/log"
	"gateway/src/mq"
	"gateway/src/pb"
	"github.com/golang/protobuf/proto"
	"time"
)

type NatsUCenter struct {
	Server   *Server
	log      *log.Logger
	NatsPool *mq.NatsPool
}

func NewNatsUCenter(server *Server) *NatsUCenter {
	return &NatsUCenter{
		Server:   server,
		log:      server.Log,
		NatsPool: server.NatsPool,
	}
}

func (self *NatsUCenter) ApplyUid(pid string) (string, error) {
	self.log.Info("ApplyUid begin pid %+v", pid)
	request, _ := proto.Marshal(&pb.ApplyUidRequest{Pid: pid})
	var response interface{}
	err := self.NatsPool.Request(constants.UCENTER_APPLY_UID_SUBJECT, string(request), &response, 3*time.Second)
	if err != nil {
		self.log.Error("applyUid subject %+v err %+v", constants.UCENTER_APPLY_UID_SUBJECT, err)
		return "", err
	}
	dataMap := response.(map[string]interface{})
	bytes := []byte(dataMap["data"].(string))
	var res pb.ApplyUidResponse
	proto.Unmarshal(bytes, &res)

	self.log.Info("ApplyUid response uid [%+v] pid [%+v] code [%+v]",
		res.GetUid(), res.GetPid(), res.GetCode())

	return res.GetUid(), nil
}
