package handler

import (
	"github.com/nats-io/go-nats"
	"strconv"
	"ucenter/src"
	"ucenter/src/constants"
	"ucenter/src/log"
	"ucenter/src/mq"
)

type HandlerMgr struct {
	server      *ucenter.Server
	NatsPool    *mq.NatsPool
	MsgFromNats chan ucenter.Closure
	log         *log.Logger
	exit        chan bool
}

func NewHandler(server *ucenter.Server) *HandlerMgr {
	return &HandlerMgr{
		server:      server,
		NatsPool:    server.NatsPool,
		MsgFromNats: make(chan ucenter.Closure, 4096),
		log:         server.Log,
		exit:        make(chan bool),
	}
}

func (self *HandlerMgr) Run() {
	defer func() {
		p := recover()
		if p != nil {
			self.log.Info("execute panic recovered and going to stop: %v", p)
		}
	}()

	self.SubscribeGetUid()

	for {
		// 优先查看exit，
		select {
		case <-self.exit:
			{
				return
			}

		default:
			// do nothing
		}
	}
}

func (self *HandlerMgr) GreateUid() string {
	uid, _ := self.server.RedisDao.IncrBy("create_uid", 1)
	return strconv.FormatInt(uid+1000, 10)
}

func (self *HandlerMgr) SubscribeGetUid() {
	// 订阅一个Nats Request 主题
	err := self.NatsPool.SubscribeForRequest(constants.UCENTER_APPLY_UID_SUBJECT, func(subj, reply string, msg interface{}) {
		self.log.Info("Nats Subscribe request subject:%+v,receive massage:%+v,reply subject:%+v", subj, msg, reply)

		natsMsg, ok := msg.(*nats.Msg)
		if ok {
			self.GetPlayerUID(reply, natsMsg)
		} else {
			self.log.Error("SubscribeGetUid Failed to convert interface{} to *nats.Msg")
		}

	})

	if err != nil {
		self.log.Error("SubscribeGetUid err %+v", err)
	}
}

func (self *HandlerMgr) Quit() {
	self.NatsPool.Unsubscribe(constants.UCENTER_APPLY_UID_SUBJECT)
	self.exit <- true
}
