package match

import (
	"encoding/json"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"match/src/domain"
	"match/src/internal"
)

type DynamicConfig struct {
	Server      *Server
	nacosConfig *domain.NacosConfig
	GameMap     map[string]*domain.GameInfo
	exit        chan bool
}

func NewDynamicConfig(server *Server) *DynamicConfig {
	config := &DynamicConfig{
		Server:      server,
		nacosConfig: server.Config.NacosConfig,
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

	self.loadAllGameData()

	_ = internal.ConfigClient.ListenConfig(vo.ConfigParam{
		DataId: self.nacosConfig.GameDataId,
		Group:  self.nacosConfig.GameGroup,
		OnChange: func(namespace, group, dataId, data string) {
			internal.GLog.Info("config changed group: %+v dataId: %v content: %+v", group, dataId, data)
			self.syncGameData(data)
		},
	})

	self.Server.WaitGroup.Add(1)
	defer func() {
		self.Server.WaitGroup.Done()
	}()

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
	gameIds, err := internal.RedisDao.SMembers("ALL_GAME_ID")
	if err != nil {
		internal.GLog.Error("load all game data error: %+v", err)
		return
	}

	for _, gameId := range gameIds {
		self.Server.SubscribeGame.SubscribeGame(gameId)
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
	internal.GLog.Info("dynamic config quit")
	self.exit <- true
}

func (self *DynamicConfig) GetGameInfo(gameId string) *domain.GameInfo {
	return self.GameMap[gameId]
}
