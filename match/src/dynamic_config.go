package match

import (
	"encoding/json"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"match/src/domain"
	"match/src/log"
	"match/src/redis"
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

	self.loadAllGameData()

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

func (self *DynamicConfig) GetAllGameId() []string {
	result := make([]string, 0)
	for key, _ := range self.GameMap {
		result = append(result, key)
	}
	return result
}

func (self *DynamicConfig) loadAllGameData() {
	//get config
	gameIds, err := self.redisDao.SMembers("ALL_GAME_ID")
	if err != nil {
		self.log.Error("load all game data error: %+v", err)
		return
	}

	for _, gameId := range gameIds {
		self.Server.SubscribeGame.SubscribeGame(gameId)
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
	self.exit <- true
}

func (self *DynamicConfig) GetGameInfo(gameId string) *domain.GameInfo {
	return self.GameMap[gameId]
}
