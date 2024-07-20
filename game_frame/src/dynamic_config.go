package game_frame

import (
	"encoding/json"
	"game_frame/src/domain"
	"game_frame/src/log"
	"game_frame/src/redis"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

type DynamicConfig struct {
	Server      *Server
	NacosConfig *NacosConfig
	redisDao    *redis.RedisDao
	log         *log.Logger
	GameInfo    *domain.GameInfo
	Client      config_client.IConfigClient
	gameChange  chan string
	exit        chan bool
}

func NewDynamicConfig(server *Server) *DynamicConfig {
	config := &DynamicConfig{
		Server:      server,
		redisDao:    server.RedisDao,
		NacosConfig: server.Config.NacosConfig,
		log:         server.Log,
		GameInfo:    nil,
		Client:      server.ConfigClient,
		gameChange:  make(chan string, 100),
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

	self.syncGameData(self.Server.Config.GameId)

	_ = self.Client.ListenConfig(vo.ConfigParam{
		DataId: self.NacosConfig.GameDataId,
		Group:  self.NacosConfig.GameGroup,
		OnChange: func(namespace, group, dataId, data string) {
			self.log.Info("config changed group: %+v dataId: %v content: %+v", group, dataId, data)
			self.gameChange <- data
		},
	})

	for {
		// 优先查看exit，
		select {
		case <-self.exit:
			{
				return
			}

		case <-self.gameChange:
			{
				gameId := <-self.gameChange

				if self.GameInfo.GameId == gameId {
					self.syncGameData(gameId)
				}
			}

		default:
			// do nothing
		}
	}
}

func (self *DynamicConfig) syncGameData(gameId string) {
	data, _ := self.redisDao.Get(gameId)
	self.log.Info("syncGameData gameId %+v data %+v", gameId, data)
	var gameInfo domain.GameInfo
	json.Unmarshal([]byte(data), &gameInfo)
	self.GameInfo = &gameInfo

}

func (self *DynamicConfig) Quit() {
	self.exit <- true
}

func (self *DynamicConfig) GetGameInfo() *domain.GameInfo {
	return self.GameInfo
}
