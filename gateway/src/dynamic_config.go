package gateway

import (
	"encoding/json"
	"gateway/src/domain"
	"gateway/src/log"
	"gateway/src/redis"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

type DynamicConfig struct {
	Server      *Server
	NacosConfig *NacosConfig
	redisDao    *redis.RedisDao
	log         *log.Logger
	GameMap     map[string]*domain.GameInfo
	Client      config_client.IConfigClient
	exit        chan bool
}

func NewDynamicConfig(server *Server) *DynamicConfig {
	config := &DynamicConfig{
		Server:      server,
		redisDao:    server.RedisDao,
		NacosConfig: server.Config.NacosConfig,
		log:         server.Log,
		GameMap:     make(map[string]*domain.GameInfo),
		Client:      server.ConfigClient,
		exit:        make(chan bool, 1),
	}

	return config
}

func (self *DynamicConfig) Run() {
	defer func() {
		p := recover()
		if p != nil {
			self.log.Info("execute panic recovered and going to stop: %v", p)
		}
	}()

	self.log.Info("dynamic config  begin ....")

	self.Server.WaitGroup.Add(1)
	defer func() {
		self.Server.WaitGroup.Done()
	}()

	self.LoadAllGameData()

	_ = self.Client.ListenConfig(vo.ConfigParam{
		DataId: self.NacosConfig.GameDataId,
		Group:  self.NacosConfig.GameGroup,
		OnChange: func(namespace, group, dataId, data string) {
			self.log.Info("config changed group: %+v dataId: %v content: %+v", group, dataId, data)
			self.syncGameData(data)
		},
	})

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
}

func (self *DynamicConfig) LoadAllGameData() {
	//get config
	gameIds, err := self.redisDao.SMembers("ALL_GAME_ID")
	if err != nil {
		self.log.Error("load all game data error: %+v", err)
		return
	}

	for _, gameId := range gameIds {
		self.syncGameData(gameId)
	}
}

func (self *DynamicConfig) syncGameData(gameId string) {
	data, _ := self.redisDao.Get(gameId)
	self.log.Info("syncGameData gameId %+v data %+v", gameId, data)
	var gameInfo domain.GameInfo
	json.Unmarshal([]byte(data), &gameInfo)
	self.GameMap[gameId] = &gameInfo

}

func (self *DynamicConfig) GetAllGame() map[string]*domain.GameInfo {
	return self.GameMap
}

func (self *DynamicConfig) Quit() {
	self.log.Info("dynamic config quit ...")
	self.exit <- true
}

func (self *DynamicConfig) GetGameInfo(gameId string) *domain.GameInfo {
	return self.GameMap[gameId]
}
