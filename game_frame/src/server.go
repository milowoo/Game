package game_frame

import (
	"game_frame/src/log"
	"game_frame/src/mq"
	"game_frame/src/redis"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"os"
	"os/signal"
	"sync"
)

var g_Server *Server

type Server struct {
	Log           *log.Logger
	Config        *GlobalConfig
	DynamicConfig *DynamicConfig
	NatsService   *NatsService
	WaitGroup     *sync.WaitGroup
	NatsPool      *mq.NatsPool
	RedisDao      *redis.RedisDao
	ConfigClient  config_client.IConfigClient
	NameClient    naming_client.INamingClient
	RoomMgr       *RoomMgr
	isQuit        bool
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

	pool, err := mq.NatsInit(config.NatsConfig.Address)
	if err != nil {
		log.Error("NewServer nat init err")
		return nil, err
	}

	g_Server.NatsPool = pool

	g_Server.ConfigClient, g_Server.NameClient, err = g_Server.InitNacos(config)
	if err != nil {
		log.Error("nacos 连接失败", err)
		return nil, err
	}

	g_Server.DynamicConfig = NewDynamicConfig(g_Server)

	g_Server.NatsService = NewNatsService(g_Server)

	g_Server.RoomMgr = NewRoomMgr(g_Server)

	g_Server.WaitGroup.Add(1) // 对应Quit中的Done

	return g_Server, nil
}

// 通知服务器退出
func (self *Server) Quit() {
	if self.isQuit {
		return
	}

	self.isQuit = true
	self.NatsService.Quit()
	self.RoomMgr.Quit()
	self.DynamicConfig.Quit()

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

	go self.NatsService.Run()
	go self.RoomMgr.Run()
	go self.DynamicConfig.Run()

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

func (server *Server) InitNacos(config *GlobalConfig) (config_client.IConfigClient, naming_client.INamingClient, error) {
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
		return nil, nil, err
	}

	// create naming client
	client, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)

	if err != nil {
		log.Error("nacos 连接失败 %+v", err)
		panic(err)
		return nil, nil, err
	}

	//Register
	param := vo.RegisterInstanceParam{
		Ip:          GetHostIp(),
		Port:        config.NacosConfig.Port,
		ServiceName: config.GameId,
		GroupName:   config.NacosConfig.GameGroup,
		ClusterName: "cluster-a",
		Weight:      10,
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
		Metadata:    map[string]string{"idc": "shanghai"},
	}

	success, err := client.RegisterInstance(param)
	if !success || err != nil {
		panic("RegisterServiceInstance failed!" + err.Error())
		return nil, nil, err
	}

	log.Info("RegisterServiceInstance,param:%+v,result:%+v ", param, success)

	return configClient, client, nil
}
