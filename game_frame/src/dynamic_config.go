package game_frame

import (
	"encoding/json"
	"game_frame/src/domain"
	"game_frame/src/internal"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

type DynamicConfig struct {
	Server   *Server
	GameInfo *domain.GameInfo
	exit     chan bool
}

func NewDynamicConfig(server *Server) *DynamicConfig {
	config := &DynamicConfig{
		Server:   server,
		GameInfo: nil,
		exit:     make(chan bool, 1),
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

	self.syncGameData(self.Server.Config.GameId)

	_ = internal.ConfigClient.ListenConfig(vo.ConfigParam{
		DataId: self.Server.Config.NacosConfig.GameDataId,
		Group:  self.Server.Config.NacosConfig.GameGroup,
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

func (self *DynamicConfig) syncGameData(gameId string) {
	data, _ := internal.RedisDao.Get(gameId)
	internal.GLog.Info("syncGameData gameId %+v data %+v", gameId, data)
	var gameInfo domain.GameInfo
	json.Unmarshal([]byte(data), &gameInfo)
	self.GameInfo = &gameInfo

}

func (self *DynamicConfig) Quit() {
	internal.GLog.Info(" dynamic config quit game")
	self.exit <- true
}

func (self *DynamicConfig) GetGameInfo() *domain.GameInfo {
	return self.GameInfo
}
