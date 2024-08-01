package ucenter

import (
	"github.com/golang/protobuf/proto"
	"os"
	"os/signal"
	"sync"
	"ucenter/src/config"
	"ucenter/src/constants"
	"ucenter/src/handler"
	"ucenter/src/internal"
	"ucenter/src/mongo"
	"ucenter/src/mq"
	"ucenter/src/pb"
	"ucenter/src/redis"
	"ucenter/src/utils"
)

var g_Server *Server

type Server struct {
	WaitGroup *sync.WaitGroup
	isQuit    bool
}

func NewServer() (*Server, error) {
	log := internal.GLog
	config, err := config.NewGlobalConfig()
	if err != nil {
		log.Error("NewServer log config err")
		return nil, err
	}

	g_Server = &Server{
		WaitGroup: &sync.WaitGroup{},
	}

	client := mongo.Connect(config.MongoConfig.Address, config.MongoConfig.Name, log)
	if client == nil {
		log.Error("NewDataSource err  ")
		return nil, nil
	}

	internal.Mongo = mongo.NewMongoDao(client, log)

	internal.RedisDao = redis.NewRedis(config.RedisConfig.Address, config.RedisConfig.MasterName, config.RedisConfig.Password)

	pool, err := mq.NatsInit(config.NatsConfig.Address)
	if err != nil {
		log.Error("NewServer nat init err")
		return nil, err
	}

	log.Info("nats address %+v success ", config.NatsConfig.Address)

	internal.NatsPool = pool

	g_Server.WaitGroup.Add(1) // 对应Quit中的Done

	return g_Server, nil
}

// 通知服务器退出
func (self *Server) Quit() {
	if self.isQuit {
		return
	}

	internal.GLog.Info("server is quit")
	self.isQuit = true

	self.WaitGroup.Done()
}

func (self *Server) Join() {
	self.WaitGroup.Wait()
}

func (self *Server) Run() {
	log := internal.GLog

	log.Info("ucenter run begin ")

	self.SubscribeGetUid()

	// wait exit
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, os.Kill)
	go func() {
		<-c
		log.Warn("exit svr by signal ...")
		self.Quit()
	}()

	self.Join()
}

func (self *Server) SubscribeGetUid() {
	log := internal.GLog
	log.Info("SubscribeGetUid subject %+v begin ... ", constants.UCENTER_APPLY_UID_SUBJECT)

	// 订阅一个Nats Request 主题
	err := internal.NatsPool.SubscribeForRequest(constants.UCENTER_APPLY_UID_SUBJECT, func(subj, reply string, msg interface{}) {
		log.Info("Nats Subscribe request subject:%+v,receive massage:%+v,reply subject:%+v", subj, msg, reply)
		req, _ := utils.ConvertInterfaceToString(msg)
		var request pb.ApplyUidRequest
		proto.Unmarshal([]byte(req), &request)
		log.Info("SubscribeGetUid request pid %+v ", request.GetPid())

		handler.GetPlayerUID(reply, &request)
	})

	if err != nil {
		log.Error("SubscribeGetUid err %+v", err)
	}
}
