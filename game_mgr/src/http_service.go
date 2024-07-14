package game_mgr

import (
	"encoding/json"
	"fmt"
	"game_mgr/src/log"
	"game_mgr/src/mongo"
	"game_mgr/src/pb"
	"game_mgr/src/redis"
	"github.com/gogo/protobuf/proto"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/nats-io/go-nats"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

// http 服务主要处理匹配回调
type HttpService struct {
	gameConfig *GlobalConfig
	httpSvr    *http.Server
	redisDao   *redis.RedisDao
	mongo      *mongo.MongoDAO
	natsConn   *nats.Conn
	iClient    config_client.IConfigClient
	log        *log.Logger
	hostIp     string
	pid        string
	StopChan   chan bool
}

type Response struct {
	Code int    `json:"code,omitempty"`
	Msg  string `json:"msg,omitempty"`
	Data string `json:"data,omitempty"`
}

// 定义一个枚举类型
type GameType int

// 使用 iota 定义枚举常量
const (
	UnKnow GameType = iota
	Single
	OneVOne
	HallAndSingle
	HallAnd1V1
)

type GameInfo struct {
	GameId    string    `json:"gameId,omitempty"`
	Type      int32     `json:"type,omitempty"`
	Name      string    `json:"name,omitempty"`
	Status    int32     `json:"status,omitempty"`
	GameTime  int32     `json:"gameTime,omitempty"`
	MatchTime int32     `json:"matchTime,omitempty"`
	Operator  string    `json:"operator,omitempty"`
	UTime     time.Time `bson:"utime" json:"utime"`
	CTime     time.Time `bson:"ctime" json:"ctime"`
}

type GmCode struct {
	Code     int32  `json:"code,omitempty"`
	GameId   string `json:"gameId,omitempty"`
	Uid      string `json:"uid,omitempty"`
	Opt      string `json:"opt,omitempty"`
	SignStr  string `json:"signStr,omitempty"`
	Operator string `json:"operator,omitempty"`
}

func NewHttpService(config *GlobalConfig, log *log.Logger) (*HttpService, error) {
	if config.Port == 0 {
		return nil, nil
	}

	httpService := &HttpService{
		gameConfig: config,
		log:        log,
		httpSvr:    nil,
		StopChan:   make(chan bool),
	}

	httpService.pid = fmt.Sprintf("%d", os.Getpid())

	redisDao := redis.NewRedis(config.RedisConfig.Address, config.RedisConfig.MasterName, config.RedisConfig.Password)
	httpService.redisDao = redisDao

	databaseDAO, err := mongo.NewDataSource(config.MongoConfig.Address, config.MongoConfig.Name, log)
	if err != nil {
		log.Error("NewDataSource err %v", err)
		return nil, err
	}

	httpService.mongo = mongo.NewMongoDao(databaseDAO, log)

	nc, err := nats.Connect(config.NatsConfig.Address)
	if nc == nil {
		log.Error("nats连接失败", err)
		return nil, err
	}
	defer nc.Close()
	httpService.natsConn = nc

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
	client, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)

	httpService.iClient = client

	hostIP := "127.0.0.1"
	if !config.LocalModel {
		hostIP = os.Getenv("K8S_HOST_IP")
	}
	httpService.hostIp = fmt.Sprintf("%v:%v", hostIP, config.Port)
	return httpService, nil
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

		self.log.Info(" uri %+v receive %+v ", r.RequestURI, string(body))
		switch r.RequestURI {
		case "/addGame":
			//新增游戏
			self.addGame(body, w)
		case "/updateGame":
			//游戏数据更新
			self.updateGame(body, w)
		case "/onOffGame":
			//游戏上下架
			self.onOffGame(body, w)
		case "/gmCode":
			//GM CODE 命令
			self.gmCode(body, w)
		default:
			httpRes := Response{Code: 500, Msg: "invalid request url", Data: ""}
			buf, _ := json.Marshal(httpRes)
			io.WriteString(w, string(buf))
		}
	})

	log.Info(fmt.Sprintf("hostIP %v ", self.hostIp))

	self.httpSvr = &http.Server{Addr: self.hostIp, Handler: handler}
	go func() {
		self.log.Info("  listen on %s", self.hostIp)
		err := self.httpSvr.ListenAndServe()
		if err != nil {
			self.log.Error("error, listenHttpServer %p", err)
		}
	}()
}

func (self *HttpService) addGame(body []byte, w http.ResponseWriter) {
	var request GameInfo
	err := json.Unmarshal(body, &request)
	if err != nil {
		self.log.Info(" addGame err %v ", err)
		httpRes := Response{Code: 500, Msg: "invalid request body", Data: ""}
		buf, _ := json.Marshal(httpRes)
		io.WriteString(w, string(buf))
		return
	}

	request.CTime = time.Now()
	request.UTime = time.Now()

	//判断是否有记录
	data, _ := self.redisDao.EXISTS(request.GameId)
	if data {
		self.log.Info(" addGame err have exist gameId %+v err %+v ", request.GameId, err)
		httpRes := Response{Code: 500, Msg: "game have exist", Data: ""}
		buf, _ := json.Marshal(httpRes)
		io.WriteString(w, string(buf))
		return
	}

	//判断是否存在这个game 信息
	_, err = self.mongo.GetGame(request.GameId)
	if err != nil {
		self.log.Info(" addGame err have exist gameId %+v err %+v ", request.GameId, err)
		httpRes := Response{Code: 500, Msg: "game have exist", Data: ""}
		buf, _ := json.Marshal(httpRes)
		io.WriteString(w, string(buf))
		return
	}

	redisBuf, _ := json.Marshal(request)
	self.redisDao.Set(request.GameId, redisBuf, 0)

	self.redisDao.SAdd("ALL_GAME_ID", request.GameId)

	self.mongo.InsertGame(&request)

	//通告 gameId 变化
	_, err = self.iClient.PublishConfig(vo.ConfigParam{
		DataId:  self.gameConfig.NacosConfig.GameDataId,
		Group:   self.gameConfig.NacosConfig.GameGroup,
		Content: request.GameId,
	})

	httpRes := Response{Code: 100, Msg: "success", Data: ""}
	buf, _ := json.Marshal(httpRes)
	io.WriteString(w, string(buf))
}

func (self *HttpService) updateGame(body []byte, w http.ResponseWriter) {
	var request GameInfo
	err := json.Unmarshal(body, &request)
	if err != nil {
		self.log.Info(" addGame err %v ", err)
		httpRes := Response{Code: 500, Msg: "invalid request body", Data: ""}
		buf, _ := json.Marshal(httpRes)
		io.WriteString(w, string(buf))
		return
	}

	//获取 game 信息
	dbData, err := self.mongo.GetGame(request.GameId)
	if err == nil {
		self.log.Info(" addGame err no exist gameId %+v err %+v ", request.GameId, err)
		httpRes := Response{Code: 500, Msg: "game no exist", Data: ""}
		buf, _ := json.Marshal(httpRes)
		io.WriteString(w, string(buf))
		return
	}

	if request.GameTime != 0 {
		dbData.GameTime = request.GameTime
	}
	if request.Type != 0 {
		dbData.Type = request.Type
	}

	if len(request.Name) > 1 {
		dbData.Name = request.Name
	}

	if request.MatchTime != 0 {
		dbData.MatchTime = request.MatchTime
	}

	dbData.UTime = time.Now()
	if len(request.Operator) > 0 {
		dbData.Operator = request.Operator
	}

	//记录 redis
	redisBuf, _ := json.Marshal(dbData)
	self.redisDao.Set(request.GameId, redisBuf, 0)

	self.mongo.SavaGame(request.GameId, &dbData)

	//通告 gameId 变化
	_, err = self.iClient.PublishConfig(vo.ConfigParam{
		DataId:  self.gameConfig.NacosConfig.GameDataId,
		Group:   self.gameConfig.NacosConfig.GameGroup,
		Content: request.GameId,
	})

	httpRes := Response{Code: 100, Msg: "success", Data: ""}
	buf, _ := json.Marshal(httpRes)
	io.WriteString(w, string(buf))

}
func (self *HttpService) onOffGame(body []byte, w http.ResponseWriter) {
	var request GameInfo
	err := json.Unmarshal(body, &request)
	if err != nil {
		self.log.Info(" onOffGame err %v ", err)
		httpRes := Response{Code: 500, Msg: "invalid request body", Data: ""}
		buf, _ := json.Marshal(httpRes)
		io.WriteString(w, string(buf))
		return
	}

	if request.Status != 1 && request.Status != 2 {
		self.log.Info(" onOffGame err %v ", err)
		httpRes := Response{Code: 500, Msg: "invalid request status", Data: ""}
		buf, _ := json.Marshal(httpRes)
		io.WriteString(w, string(buf))
		return
	}

	dbData, err := self.mongo.GetGame(request.GameId)
	if err == nil {
		self.log.Info(" addGame err no exist gameId %+v err %+v ", request.GameId, err)
		httpRes := Response{Code: 500, Msg: "game no exist", Data: ""}
		buf, _ := json.Marshal(httpRes)
		io.WriteString(w, string(buf))
		return
	}

	dbData.Status = request.Status
	dbData.UTime = time.Now()
	if len(request.Operator) > 0 {
		dbData.Operator = request.Operator
	}

	redisBuf, _ := json.Marshal(dbData)
	self.redisDao.Set(request.GameId, redisBuf, 0)

	self.mongo.SavaGame(request.GameId, &dbData)

	//通告 gameId 变化
	_, err = self.iClient.PublishConfig(vo.ConfigParam{
		DataId:  self.gameConfig.NacosConfig.GameDataId,
		Group:   self.gameConfig.NacosConfig.GameGroup,
		Content: request.GameId,
	})

	httpRes := Response{Code: 100, Msg: "success", Data: ""}
	buf, _ := json.Marshal(httpRes)
	io.WriteString(w, string(buf))
}
func (self *HttpService) gmCode(body []byte, w http.ResponseWriter) {
	var request GmCode
	err := json.Unmarshal(body, &request)
	if err != nil {
		self.log.Info(" gmCode err %v ", err)
		httpRes := Response{Code: 500, Msg: "invalid request body", Data: ""}
		buf, _ := json.Marshal(httpRes)
		io.WriteString(w, string(buf))
		return
	}

	//验签

	builder, _ := proto.Marshal(&pb.GmCodeRequest{
		Code:     request.Code,
		GameId:   request.GameId,
		Uid:      request.Uid,
		Opt:      request.Opt,
		Operator: request.Operator,
	})

	subject := "GMCODE." + request.GameId

	_, err = self.natsConn.Request(subject, builder, 3*time.Second)
	if err != nil {
		self.log.Error("applyUid err %+v", err)
		httpRes := Response{Code: 500, Msg: "request game err", Data: ""}
		buf, _ := json.Marshal(httpRes)
		io.WriteString(w, string(buf))
		return
	}

	httpRes := Response{Code: 100, Msg: "success", Data: ""}
	buf, _ := json.Marshal(httpRes)
	io.WriteString(w, string(buf))
	return
}

func (self *HttpService) OnQuit() {
	self.StopChan <- true
}
