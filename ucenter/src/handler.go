package ucenter

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"reflect"
	"strconv"
	"ucenter/src/constants"
	"ucenter/src/log"
	"ucenter/src/mq"
	"ucenter/src/pb"
)

type HandlerMgr struct {
	server      *Server
	NatsPool    *mq.NatsPool
	MsgFromNats chan Closure
	log         *log.Logger
	exit        chan bool
}

func NewHandler(server *Server) *HandlerMgr {
	return &HandlerMgr{
		server:      server,
		NatsPool:    server.NatsPool,
		MsgFromNats: make(chan Closure, 10*1024),
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

	self.server.WaitGroup.Add(1)
	defer func() {
		self.server.WaitGroup.Done()
	}()

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

	self.log.Info("Run quit ...")
}

func (self *HandlerMgr) CreateUid() string {
	uid, _ := self.server.RedisDao.IncrBy("create_uid", 1)
	return strconv.FormatInt(uid+1000, 10)
}

func (self *HandlerMgr) SubscribeGetUid() {
	self.log.Info("SubscribeGetUid subject %+v begin ... ", constants.UCENTER_APPLY_UID_SUBJECT)

	// 订阅一个Nats Request 主题
	err := self.NatsPool.SubscribeForRequest(constants.UCENTER_APPLY_UID_SUBJECT, func(subj, reply string, msg interface{}) {
		self.log.Info("Nats Subscribe request subject:%+v,receive massage:%+v,reply subject:%+v", subj, msg, reply)
		dataType := reflect.TypeOf(msg)
		self.log.Info("SubscribeGetUid req type %+v ", dataType)
		req, _ := ConvertInterfaceToString(msg)
		var request pb.ApplyUidRequest
		proto.Unmarshal([]byte(req), &request)
		self.log.Info("SubscribeGetUid request pid %+v ", request.GetPid())

		self.GetPlayerUID(reply, &request)

		//response := &pb.ApplyUidResponse{
		//	Code: constants.CODE_SUCCESS,
		//	Uid:  "223",
		//	Pid:  request.Pid,
		//}
		//
		//resByte, _ := proto.Marshal(response)
		//
		//self.NatsPool.Publish(reply, map[string]interface{}{"res": "ok", "data": string(resByte)})

	})

	if err != nil {
		self.log.Error("SubscribeGetUid err %+v", err)
	}
}

func ConvertInterfaceToString(data interface{}) (string, error) {
	// 使用 reflect 包检查 data 是否为 string 类型
	if reflect.TypeOf(data).Kind() != reflect.String {
		return "", fmt.Errorf("expected a string, got %T", data)
	}

	// 如果是 string 类型，返回其数据
	return data.(string), nil
}

func (self *HandlerMgr) Quit() {
	self.log.Info("ucenter quit")
	self.NatsPool.Unsubscribe(constants.UCENTER_APPLY_UID_SUBJECT)
	self.exit <- true
}
