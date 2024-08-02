package gateway

import (
	"encoding/json"
	"gateway/src/domain"
	"gateway/src/internal"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

type DynamicConfig struct {
	Server      *Server
	NacosConfig *domain.NacosConfig
	GameMap     map[string]*domain.GameInfo
	exit        chan bool
}

func NewDynamicConfig(server *Server) *DynamicConfig {
	config := &DynamicConfig{
		Server:      server,
		NacosConfig: server.Config.NacosConfig,
		GameMap:     make(map[string]*domain.GameInfo),
		exit:        make(chan bool, 1),
	}

	return config
}

func (self *DynamicConfig) Run() {
	defer func() {
		p := recover()
		if p != nil {
			internal.GLog.Info("execute panic recovered and going to stop: %v", p)
		}
	}()

	internal.GLog.Info("dynamic config  begin ....")

	self.Server.WaitGroup.Add(1)
	defer func() {
		self.Server.WaitGroup.Done()
	}()

	self.LoadAllGameData()

	_ = internal.ConfigClient.ListenConfig(vo.ConfigParam{
		DataId: self.NacosConfig.GameDataId,
		Group:  self.NacosConfig.GameGroup,
		OnChange: func(namespace, group, dataId, data string) {
			internal.GLog.Info("config changed group: %+v dataId: %v content: %+v", group, dataId, data)
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
	gameIds, err := internal.RedisDao.SMembers("ALL_GAME_ID")
	if err != nil {
		internal.GLog.Error("load all game data error: %+v", err)
		return
	}

	for _, gameId := range gameIds {
		self.syncGameData(gameId)
	}
}

func (self *DynamicConfig) syncGameData(gameId string) {
	data, _ := internal.RedisDao.Get(gameId)
	internal.GLog.Info("syncGameData gameId %+v data %+v", gameId, data)
	var gameInfo domain.GameInfo
	json.Unmarshal([]byte(data), &gameInfo)
	self.GameMap[gameId] = &gameInfo

}

func (self *DynamicConfig) GetAllGame() map[string]*domain.GameInfo {
	return self.GameMap
}

func (self *DynamicConfig) Quit() {
	internal.GLog.Info("dynamic config quit ...")
	self.exit <- true
}

func (self *DynamicConfig) GetGameInfo(gameId string) *domain.GameInfo {
	return self.GameMap[gameId]
}
