package ucenter

import (
	"os"
	"os/signal"
	"sync"
	"ucenter/src/log"
	"ucenter/src/mongo"
	"ucenter/src/mq"
	"ucenter/src/redis"
)

var g_Server *Server

type Server struct {
	Log       *log.Logger
	Config    *GlobalConfig
	WaitGroup *sync.WaitGroup
	NatsPool  *mq.NatsPool
	RedisDao  *redis.RedisDao
	MongoDao  *mongo.MongoDAO
	Handler   *HandlerMgr
	isQuit    bool
}

func NewServer(log *log.Logger) (*Server, error) {
	config, err := NewGlobalConfig(log)
	if err != nil {
		log.Error("NewServer log config err")
		return nil, err
	}

	g_Server = &Server{
		Config:    config,
		WaitGroup: &sync.WaitGroup{},
		Log:       log,
	}

	client := mongo.Connect(config.MongoConfig.Address, config.MongoConfig.Name, log)
	if client == nil {
		log.Error("NewDataSource err  ")
		return nil, nil
	}

	g_Server.MongoDao = mongo.NewMongoDao(client, log)

	redisDao := redis.NewRedis(config.RedisConfig.Address, config.RedisConfig.MasterName, config.RedisConfig.Password)
	g_Server.RedisDao = redisDao

	pool, err := mq.NatsInit(config.NatsConfig.Address)
	if err != nil {
		log.Error("NewServer nat init err")
		return nil, err
	}

	log.Info("nats address %+v success ", config.NatsConfig.Address)

	g_Server.NatsPool = pool

	handlerMgr := NewHandler(g_Server)
	g_Server.Handler = handlerMgr

	g_Server.WaitGroup.Add(1) // 对应Quit中的Done

	return g_Server, nil
}

// 通知服务器退出
func (self *Server) Quit() {
	if self.isQuit {
		return
	}

	self.Log.Info("server is quit")
	self.isQuit = true
	self.Handler.Quit()

	self.WaitGroup.Done()
}

func (self *Server) Join() {
	self.WaitGroup.Wait()
}

func (self *Server) Run() {
	if self.Config == nil {
		self.Log.Error("Config is nil ...")
		self.Quit()
	}

	self.Log.Info("ucenter run begin ")

	go self.Handler.Run()

	// wait exit
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, os.Kill)
	go func() {
		<-c
		self.Log.Warn("exit svr by signal ...")
		self.Quit()
	}()

	self.Join()
}
