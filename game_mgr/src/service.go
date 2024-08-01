package game_mgr

import (
	"encoding/json"
	"fmt"
	"game_mgr/src/config"
	"game_mgr/src/constants"
	"game_mgr/src/domain"
	"game_mgr/src/handler"
	"game_mgr/src/internal"
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
	GameConfig *config.GlobalConfig
	httpSvr    *http.Server
	hostIp     string
	pid        string
	StopChan   chan bool
}

func NewHttpService(globalConfig *config.GlobalConfig) (*HttpService, error) {
	if globalConfig.Port == 0 {
		return nil, nil
	}

	httpService := &HttpService{
		httpSvr:    nil,
		GameConfig: globalConfig,
		StopChan:   make(chan bool),
	}

	httpService.pid = fmt.Sprintf("%d", os.Getpid())

	internal.RedisDao = redis.NewRedis(globalConfig.RedisConfig.Address, globalConfig.RedisConfig.MasterName, globalConfig.RedisConfig.Password)

	client := mongo.Connect(globalConfig.MongoConfig.Address, globalConfig.MongoConfig.Name, internal.GLog)
	if client == nil {
		internal.GLog.Error("NewDataSource err  ")
		return nil, nil
	}

	internal.Mongo = mongo.NewMongoDao(client, internal.GLog)

	natsPool, err := mq.NatsInit(globalConfig.NatsConfig.Address)
	if err != nil {
		internal.GLog.Error("nats 连接失败 %+v", err)
		return nil, err
	}

	internal.NatsPool = natsPool

	internal.IClient, err = httpService.InitNacos(globalConfig)
	if err != nil {
		internal.GLog.Error("nats连接失败", err)
		return nil, err
	}

	hostIP := "127.0.0.1"
	if !globalConfig.LocalModel {
		hostIP = os.Getenv("K8S_HOST_IP")
	}
	httpService.hostIp = fmt.Sprintf("%+v:%+v", hostIP, globalConfig.Port)
	internal.GLog.Info("hostIp %+v Port %+v", hostIP, globalConfig.Port)
	return httpService, nil
}

func (server *HttpService) InitNacos(globalConfig *config.GlobalConfig) (config_client.IConfigClient, error) {
	sc := []constant.ServerConfig{
		*constant.NewServerConfig(globalConfig.NacosConfig.Ip, globalConfig.NacosConfig.Port, constant.WithContextPath("/nacos")),
	}

	//create ClientConfig
	cc := *constant.NewClientConfig(
		constant.WithNamespaceId(globalConfig.NacosConfig.SpaceId),
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

		internal.GLog.Info(" uri %+v receive %+v ", r.RequestURI, string(body))
		switch r.RequestURI {
		case "/addGame":
			//新增游戏
			handler.AddGameHandler(self.GameConfig.NacosConfig, body, w)
		case "/updateGame":
			//游戏数据更新
			handler.UpdateGameHandler(self.GameConfig.NacosConfig, body, w)
		case "/onOffGame":
			//游戏上下架
			handler.OnOffGameHandler(self.GameConfig.NacosConfig, body, w)
		case "/greenBlueDeployment":
			//游戏蓝绿发布
			handler.GreenBlueDeploymentHandler(self.GameConfig.NacosConfig, body, w)
		case "/addAiInfo":
			//增加AI信息
			handler.AddAiHandler(body, w)
		case "/addColoredUidList":
			//增加活动染色用户
			handler.AddColoredUidHandler(body, w)
		case "/deleteColoredUidList":
			//删除活动染色用户
			handler.DeleteColoredUidHandler(body, w)
		case "/kickColorUidList":
			//剔除染色用户
			handler.KickColoredUidHandler(body, w)
		case "/gmCode":
			//GM CODE 命令
			handler.GmCodeHandler(body, w)
		default:
			httpRes := domain.Response{Code: constants.INVALID_URL, Msg: "invalid request url", Data: ""}
			buf, _ := json.Marshal(httpRes)
			io.WriteString(w, string(buf))
		}
	})

	internal.GLog.Info(fmt.Sprintf("hostIP %v ", self.hostIp))

	self.httpSvr = &http.Server{Addr: self.hostIp, Handler: handler}
	go func() {
		internal.GLog.Info("  listen on %s", self.hostIp)
		err := self.httpSvr.ListenAndServe()
		if err != nil {
			internal.GLog.Error("error, listenHttpServer %p", err)
		}
	}()
}

func (self *HttpService) Quit() {
	self.StopChan <- true
}
