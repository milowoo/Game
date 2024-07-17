package gateway

import (
	"context"
	"fmt"
	"gateway/src/log"
	"gateway/src/mq"
	"gateway/src/redis"
	"github.com/gorilla/websocket"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"
)

var g_Server *Server

type Server struct {
	Log           *log.Logger
	Config        *GlobalConfig
	WaitGroup     *sync.WaitGroup
	AgentMgr      *AgentMgr
	DynamicConfig *DynamicConfig
	NatsPool      *mq.NatsPool
	MatchMgr      *NatsMatch
	UCenterMgr    *NatsUCenter
	RedisDao      *redis.RedisDao
	ConfigClient  config_client.IConfigClient

	svr    *http.Server
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
		svr:       nil,
		Log:       log,
	}

	agentMgr := NewAgentMgr(g_Server)
	g_Server.AgentMgr = agentMgr
	natsMatch := NewNatsMatch(g_Server)
	g_Server.MatchMgr = natsMatch
	natsUCenter := NewNatsUCenter(g_Server)
	g_Server.UCenterMgr = natsUCenter
	redisDao := redis.NewRedis(config.RedisConfig.Address, config.RedisConfig.MasterName, config.RedisConfig.Password)
	g_Server.RedisDao = redisDao
	natsPool, err := mq.NatsInit(config.NatsConfig.Address)
	if err != nil {
		log.Error("nats 连接失败 %+v", err)
		return nil, err
	}

	g_Server.NatsPool = natsPool

	g_Server.ConfigClient, err = g_Server.InitNacos(config)
	if err != nil {
		log.Error("nacos 连接失败 %+v", err)
		return nil, err
	}

	g_Server.DynamicConfig = NewDynamicConfig(g_Server)

	g_Server.WaitGroup.Add(1) // 对应Quit中的Done

	return g_Server, nil
}

var UPGRADER = websocket.Upgrader{
	ReadBufferSize:  10 * 1024,
	WriteBufferSize: 10 * 1024,
}

// 通知服务器退出

func (self *Server) Quit() {
	if self.isQuit {
		return
	}

	self.isQuit = true

	// 关闭agent监听，不再允许进入新连接
	self.StopWebsocketServer()

	// 通知所有Quit强制存储，并退出；查询中的agent直接通知服务器关闭
	self.AgentMgr.Quit()

	self.DynamicConfig.Quit()

	self.WaitGroup.Done()
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

func (self *Server) StopWebsocketServer() {
	if self.svr != nil {
		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
		self.svr.Shutdown(ctx)
	}
}

func (self *Server) Join() {
	self.WaitGroup.Wait()
}

func (self *Server) ListenAndServe() {
	if self.svr != nil {
		self.Log.Error("error: ListenAndServe() repeatedly")
		return
	}

	UPGRADER.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	clientConfig := self.Config.Client
	addr := fmt.Sprintf("%s:%d", clientConfig.ListenIp, clientConfig.ListenPort)
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rawConn, err := UPGRADER.Upgrade(w, r, nil)
		if err != nil {
			self.Log.Error("UPGRADER.Upgrade(): %+v", err)
			return
		}

		conn := NewAgent(rawConn, self.AgentMgr, r.URL.Query())
		go conn.Run()
	})
	self.svr = &http.Server{Addr: addr, Handler: handler}
	go func() {
		self.WaitGroup.Add(1) // 对应本携程
		defer self.WaitGroup.Done()

		self.Log.Info("gateway listen on %s", addr)

		err := self.svr.ListenAndServe()
		if err != nil {
			self.Log.Error("error, ListenAndServe: %+v", err)
		}

		self.svr = nil
	}()
}

func (self *Server) Run() {
	if self.Config == nil {
		self.Log.Error("Config is nil ...")
		self.Quit()
	}

	go self.AgentMgr.Run()

	go self.DynamicConfig.Run()

	go self.MatchMgr.Run()

	go self.AgentMgr.Run()

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

func (self *Server) MatchRequest(gameId string, uid string, score int32, opt string) {
	RunOnMatch(self.MatchMgr.MsgFromServer, self.MatchMgr, func(match *NatsMatch) {
		match.MatchRequest(gameId, uid, score, opt)
	})
}

func (self *Server) CancelMatchRequest(gameId string, uid string) {
	RunOnMatch(self.MatchMgr.MsgFromServer, self.MatchMgr, func(match *NatsMatch) {
		match.CancelMatchRequest(gameId, uid)
	})
}
