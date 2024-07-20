package game_mgr

import (
	"encoding/json"
	"fmt"
	"game_mgr/src/constants"
	"game_mgr/src/domain"
	"game_mgr/src/handler"
	"game_mgr/src/log"
	"game_mgr/src/mongo"
	"game_mgr/src/mq"
	"game_mgr/src/redis"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

// http 服务主要处理匹配回调
type HttpService struct {
	GameConfig *GlobalConfig
	httpSvr    *http.Server
	RedisDao   *redis.RedisDao
	Mongo      *mongo.MongoDAO
	NatsPool   *mq.NatsPool
	IClient    config_client.IConfigClient
	Log        *log.Logger
	hostIp     string
	pid        string
	StopChan   chan bool
}

func NewHttpService(config *GlobalConfig, log *log.Logger) (*HttpService, error) {
	if config.Port == 0 {
		return nil, nil
	}

	httpService := &HttpService{
		GameConfig: config,
		Log:        log,
		httpSvr:    nil,
		StopChan:   make(chan bool),
	}

	httpService.pid = fmt.Sprintf("%d", os.Getpid())

	redisDao := redis.NewRedis(config.RedisConfig.Address, config.RedisConfig.MasterName, config.RedisConfig.Password)
	httpService.RedisDao = redisDao

	databaseDAO, err := mongo.NewDataSource(config.MongoConfig.Address, config.MongoConfig.Name, log)
	if err != nil {
		log.Error("NewDataSource err %v", err)
		return nil, err
	}

	httpService.Mongo = mongo.NewMongoDao(databaseDAO, log)

	natsPool, err := mq.NatsInit(config.NatsConfig.Address)
	if err != nil {
		log.Error("nats 连接失败 %+v", err)
		return nil, err
	}

	httpService.NatsPool = natsPool

	httpService.IClient, err = httpService.InitNacos(config)
	if err != nil {
		log.Error("nats连接失败", err)
		return nil, err
	}

	hostIP := "127.0.0.1"
	if !config.LocalModel {
		hostIP = os.Getenv("K8S_HOST_IP")
	}
	httpService.hostIp = fmt.Sprintf("%v:%v", hostIP, config.Port)
	return httpService, nil
}

func (server *HttpService) InitNacos(config *GlobalConfig) (config_client.IConfigClient, error) {
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

func (self *HttpService) Run() {
	self.listenHttpServer()
	for {
		select {
		case <-self.StopChan:
			fmt.Println("http service received stop signal")
			return
		}
	}
}

func (self *HttpService) listenHttpServer() {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body)
		r.Body.Close()

		self.Log.Info(" uri %+v receive %+v ", r.RequestURI, string(body))
		switch r.RequestURI {
		case "/addGame":
			//新增游戏
			handler.AddGameHandler(self, body, w)
		case "/updateGame":
			//游戏数据更新
			handler.UpdateGameHanler(self, body, w)
		case "/onOffGame":
			//游戏上下架
			handler.OnOffGameHandler(self, body, w)
		case "/greenBlueDeployment":
			//游戏蓝绿发布
			handler.GreenBlueDeploymentHandler(self, body, w)
		case "addAiInfo":
			//增加AI信息
			handler.AddAiHandler(self, body, w)
		case "addColoredUidList":
			//增加活动染色用户
			handler.AddAiHandler(self, body, w)
		case "deleteColoredUidList":
			//删除活动染色用户
			handler.DeleteColoredUidHandler(self, body, w)
		case "kickColorUidList":
			//剔除染色用户
			handler.KickColoredUidHandler(self, body, w)
		case "/gmCode":
			//GM CODE 命令
			handler.GmCodeHandler(self, body, w)
		default:
			httpRes := domain.Response{Code: constants.INVALID_URL, Msg: "invalid request url", Data: ""}
			buf, _ := json.Marshal(httpRes)
			io.WriteString(w, string(buf))
		}
	})

	log.Info(fmt.Sprintf("hostIP %v ", self.hostIp))

	self.httpSvr = &http.Server{Addr: self.hostIp, Handler: handler}
	go func() {
		self.Log.Info("  listen on %s", self.hostIp)
		err := self.httpSvr.ListenAndServe()
		if err != nil {
			self.Log.Error("error, listenHttpServer %p", err)
		}
	}()
}

func (self *HttpService) Quit() {
	self.StopChan <- true
}
