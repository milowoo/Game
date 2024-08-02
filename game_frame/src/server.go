package game_frame

import (
	"game_frame/src/config"
	"game_frame/src/internal"
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
	Config        *config.GlobalConfig
	DynamicConfig *DynamicConfig
	NatsService   *NatsService
	WaitGroup     *sync.WaitGroup
	RoomMgr       *RoomMgr
	HostIp        string
	isQuit        bool
}

func NewServer() (*Server, error) {
	config, err := config.NewGlobalConfig()
	if err != nil {
		internal.GLog.Error("NewServer log config err")
		return nil, err
	}

	g_Server = &Server{
		Config:    config,
		HostIp:    GetHostIp(),
		WaitGroup: &sync.WaitGroup{},
	}

	redisDao := redis.NewRedis(config.RedisConfig.Address, config.RedisConfig.MasterName, config.RedisConfig.Password)
	internal.RedisDao = redisDao

	pool, err := mq.NatsInit(config.NatsConfig.Address)
	if err != nil {
		internal.GLog.Error("NewServer nat init err")
		return nil, err
	}

	internal.NatsPool = pool
	internal.ConfigClient, internal.NameClient, err = g_Server.InitNacos(config)
	if err != nil {
		internal.GLog.Error("nacos 连接失败", err)
		return nil, err
	}

	g_Server.DynamicConfig = NewDynamicConfig(g_Server)

	g_Server.NatsService = NewNatsService(g_Server)

	g_Server.RoomMgr = NewRoomMgr(g_Server)

	g_Server.registerServiceInstance(vo.RegisterInstanceParam{
		Ip:          g_Server.HostIp,
		Port:        config.Port,
		ServiceName: config.GameId,
		GroupName:   "group-a",
		ClusterName: "cluster-a",
		Weight:      10,
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
		Metadata:    map[string]string{"idc": "shanghai"},
	})

	g_Server.WaitGroup.Add(1) // 对应Quit中的Done

	return g_Server, nil
}

// 通知服务器退出
func (self *Server) Quit() {
	if self.isQuit {
		return
	}

	self.isQuit = true

	if self.NatsService != nil {
		self.NatsService.Quit()
	}

	if self.RoomMgr != nil {
		self.RoomMgr.Quit()
	}

	if self.DynamicConfig != nil {
		self.DynamicConfig.Quit()
	}

	self.deRegisterServiceInstance(vo.DeregisterInstanceParam{
		Ip:          self.HostIp,
		Port:        self.Config.Port,
		ServiceName: self.Config.GameId,
		GroupName:   "group-a",
		Cluster:     "cluster-a",
		Ephemeral:   true, //it must be true
	})

	self.WaitGroup.Done()
}

func (self *Server) Join() {
	self.WaitGroup.Wait()
}

func (self *Server) Run() {
	if self.Config == nil {
		internal.GLog.Error("Config is nil ...")
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
		internal.GLog.Warn("exit svr by signal ...")
		self.Quit()
	}()

	self.Join()
}

func (server *Server) InitNacos(config *config.GlobalConfig) (config_client.IConfigClient, naming_client.INamingClient, error) {
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
		internal.GLog.Error("nacos 连接失败 %+v", err)
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
		internal.GLog.Error("nacos 连接失败 %+v", err)
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

	internal.GLog.Info("RegisterServiceInstance,param:%+v,result:%+v ", param, success)

	return configClient, client, nil
}

func (self *Server) registerServiceInstance(param vo.RegisterInstanceParam) error {
	success, err := internal.NameClient.RegisterInstance(param)
	if !success || err != nil {
		internal.GLog.Error("RegisterServiceInstance,failed %+v \n\n", err)
		panic("RegisterServiceInstance failed!" + err.Error())
		return err
	}
	internal.GLog.Info("RegisterServiceInstance,param:%+v,result:%+v \n\n", param, success)
	return nil
}

func (self *Server) deRegisterServiceInstance(param vo.DeregisterInstanceParam) error {
	success, err := internal.NameClient.DeregisterInstance(param)
	if !success || err != nil {
		internal.GLog.Error("DeRegisterServiceInstance  %+v \n\n", err)
		panic("DeRegisterServiceInstance failed!" + err.Error())
		return err
	}
	internal.GLog.Info("DeRegisterServiceInstance,param:%+v,result:%+v \n\n", param, success)
	return nil
}
