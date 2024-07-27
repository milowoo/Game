package gateway

import (
	"errors"
	"gateway/src/constants"
	"gateway/src/log"
	"gateway/src/mq"
	"gateway/src/pb"
	"github.com/gogo/protobuf/proto"
	"github.com/nats-io/go-nats"
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
	err := self.NatsPool.Request(constants.UCENTER_APPLY_UID_SUBJECT, request, &response, 3*time.Second)
	if err != nil {
		self.log.Error("applyUid subject %+v err %+v", constants.UCENTER_APPLY_UID_SUBJECT, err)
		return "", err
	}

	natsMsg, ok := response.(*nats.Msg)
	if !ok {
		self.log.Error("Failed to convert interface{} to *nats.Msg pid %+v", pid)
		return "", err
	}

	var applyResponse pb.ApplyUidResponse
	_ = proto.Unmarshal(natsMsg.Data, &applyResponse)
	if applyResponse.Code != 0 {
		self.log.Error("applyUid err %+v %+v", applyResponse.Code, applyResponse.GetMsg())
		return "", errors.New(applyResponse.GetMsg())
	}

	return applyResponse.GetUid(), nil
}
