package match

import (
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"match/src/log"
	"match/src/mq"
	"match/src/redis"
	"os"
	"os/signal"
	"sync"
)

var g_Server *Server

type Server struct {
	Log           *log.Logger
	DynamicConfig *DynamicConfig
	SubscribeGame *SubscribeGame
	NatsService   *NatsService
	MatchMgr      map[string]*GameMatch
	Config        *GlobalConfig
	WaitGroup     *sync.WaitGroup
	NatsPool      *mq.NatsPool
	ConfigClient  config_client.IConfigClient
	NameClient    naming_client.INamingClient
	RedisDao      *redis.RedisDao

	isQuit bool
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

	redisDao := redis.NewRedis(config.RedisConfig.Address, config.RedisConfig.MasterName, config.RedisConfig.Password)
	g_Server.RedisDao = redisDao

	g_Server.NatsPool, err = mq.NatsInit(config.NatsConfig.Address)
	if err != nil {
		log.Error("NewServer nat init err")
		return nil, err
	}

	g_Server.ConfigClient, err = g_Server.InitNacos(config)
	if err != nil {
		log.Error("nacos 连接失败", err)
		return nil, err
	}

	dynamicConfig := NewDynamicConfig(g_Server)
	g_Server.DynamicConfig = dynamicConfig

	g_Server.NatsService = NewNatsService(g_Server)

	g_Server.WaitGroup.Add(1) // 对应Quit中的Done

	return g_Server, nil
}

func (server *Server) InitNacos(config *GlobalConfig) (config_client.IConfigClient, error) {
	sc := []constant.ServerConfig{
		*constant.NewServerConfig(config.NacosConfig.Ip, config.NacosConfig.Port, constant.WithContextPath("/nacos")),
	}

	//create ClientConfig
	cc := *constant.NewClientConfig(
		constant.WithNamespaceId(config.NacosConfig.SpaceId),
		constant.WithTimeoutMs(5000),
		constant.WithNotLoadCacheAtStart(true),
		constant.WithLogDir("/tmp/nacos/log"),
		constant.WithCacheDir("/tmp/nacos/cache"),
		constant.WithLogLevel("debug"),
	)

	// create config client
	configClient, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)

	if err != nil {
		log.Error("nacos 连接失败 %+v", err)
		panic(err)
		return nil, err
	}

	return configClient, nil
}

// 通知服务器退出

func (self *Server) Quit() {
	if self.isQuit {
		return
	}

	self.NatsService.Quit()

	self.DynamicConfig.Quit()
	for _, gameMath := range self.MatchMgr {
		gameMath.Quit()
	}
	self.isQuit = true

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

	go self.DynamicConfig.Run()

	gameMap := self.DynamicConfig.GetAllGame()
	for _, game := range gameMap {
		gameMatch := NewGameMatch(self, game)
		go gameMatch.Run()
		self.MatchMgr[game.GameId] = gameMatch
	}

	go self.NatsService.Run()

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
