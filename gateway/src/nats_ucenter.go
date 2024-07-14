package gateway

import (
	"errors"
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
	subject  string
	NatsPool *mq.NatsPool
}

func NewNatsUCenter(server *Server) *NatsUCenter {
	return &NatsUCenter{
		Server:   server,
		log:      server.Log,
		subject:  "ucenter.apply.uid",
		NatsPool: server.NatsPool,
	}
}

func (self *NatsUCenter) ApplyUid(app string) (string, error) {
	request, _ := proto.Marshal(&pb.ApplyUidRequest{Data: app})
	var response interface{}
	err := self.NatsPool.Request(self.subject, request, &response, 3*time.Second)
	if err != nil {
		self.log.Error("applyUid err %+v", err)
		return "", err
	}

	natsMsg, ok := response.(*nats.Msg)
	if !ok {
		self.log.Error("Failed to convert interface{} to *nats.Msg app %+v", app)
	}

	var applyResponse pb.ApplyUidResponse
	_ = proto.Unmarshal(natsMsg.Data, &applyResponse)
	if applyResponse.Code != 0 {
		self.log.Error("applyUid err %+v %+v", applyResponse.Code, applyResponse.GetMsg())
		return "", errors.New(applyResponse.GetMsg())
	}
	return applyResponse.GetUid(), nil
}
