package match

/**
监听服务 down 的处理
*/

import (
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/model"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"match/src/log"
)

type SubscribeGame struct {
	Server          *Server
	NacosConfig     *NacosConfig
	log             *log.Logger
	SubGameMap      map[string]*vo.SubscribeParam
	Client          naming_client.INamingClient
	gameChange      chan string
	subGameInstance chan string
	exit            chan bool
	isQuit          bool
}

func NewSubscribeGame(server *Server) *SubscribeGame {
	game := &SubscribeGame{
		Server:          server,
		NacosConfig:     server.Config.NacosConfig,
		log:             server.Log,
		Client:          server.NameClient,
		SubGameMap:      make(map[string]*vo.SubscribeParam),
		gameChange:      make(chan string, 100),
		subGameInstance: make(chan string, 100),
		exit:            make(chan bool, 1),
		isQuit:          false,
	}

	return game
}

func (self *SubscribeGame) Run() {
	defer func() {
		p := recover()
		if p != nil {
			self.log.Info("execute panic recovered and going to stop: %v", p)
		}
	}()

	gameSet := self.Server.DynamicConfig.GetAllGameId()
	for _, gameId := range gameSet {
		self.SubscribeGame(gameId)
	}

	for {
		// 优先查看exit，
		select {
		case <-self.exit:
			{
				self.isQuit = true
				return
			}

		case <-self.gameChange:
			{
				gameId := <-self.gameChange
				if len(gameId) > 1 {
					self.SubscribeGame(gameId)
				}
			}
		case <-self.subGameInstance:
			{
				services := <-self.subGameInstance
				if len(services) > 1 {
					self.procGameDown(services)
				}
			}

		default:
			// do nothing
		}
	}
}

func (self *SubscribeGame) SubscribeGame(gameId string) {
	if _, ok := self.SubGameMap[gameId]; ok {
		return
	}
	subscribeParam := &vo.SubscribeParam{
		ServiceName: gameId,
		GroupName:   self.NacosConfig.GameGroup,
		SubscribeCallback: func(services []model.Instance, err error) {
			service := ToJsonString(services)
			self.log.Info("callback return services %+v", service)
			self.gameChange <- service
		},
	}

	self.Client.Subscribe(subscribeParam)
	self.SubGameMap[gameId] = subscribeParam
}

func (self *SubscribeGame) procGameDown(data string) {
	services, err := JsonToServices(data)
	if err != nil {
		self.log.Error("failed to unmarshal json string:%+v err:%+v", data, err)
	}

	for _, service := range *services {
		//服务不可使用， 则需要
		self.procGameHall(service)
	}
}

func (self *SubscribeGame) procGameHall(service model.Service) {
	//取出所有的 down 处理的所有房间， 重新分配
	self.log.Info("procGameRoom service %+v", service)
	for _, instance := range service.Hosts {
		self.log.Info("procGameRoom host %+v", instance)
		if !instance.Enable {
			self.log.Error("warn service %+v", instance)
		}
	}
}

func (self *SubscribeGame) Quit() {
	if self.isQuit {
		return
	}
	for _, param := range self.SubGameMap {
		self.Client.Unsubscribe(param)
	}

	self.exit <- true
}
