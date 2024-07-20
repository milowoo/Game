package gateway

import (
	"gateway/src/log"
	"gateway/src/pb"
	"time"
)

type AgentMgr struct {
	Server       *Server
	Log          *log.Logger
	uid2Agent    map[string]*Agent
	MsgFromAgent chan Closure
	frameTimer   *time.Ticker
	frameID      int
}

func NewAgentMgr(server *Server) *AgentMgr {
	return &AgentMgr{
		Server:       server,
		Log:          server.Log,
		uid2Agent:    make(map[string]*Agent),
		MsgFromAgent: make(chan Closure, 16*1024),
		frameID:      0,
	}
}

func (self *AgentMgr) Run() {
	defer func() {
		p := recover()
		if p != nil {
			self.Log.Info("execute panic recovered and going to stop: %v", p)
		}
	}()

	self.Server.WaitGroup.Add(1)
	defer func() {
		self.Server.WaitGroup.Done()
	}()

ALL:
	for {
		select {
		case c, ok := <-self.MsgFromAgent:
			if !ok {
				break ALL
			}
			SafeRunClosure(self, c)
		case <-self.frameTimer.C:
			SafeRunClosure(self, func() {
				self.Frame()
			})
		}
	}

	self.OnQuit()
}

func (self *AgentMgr) Frame() {
	self.SetNextFrameId()

}

func (self *AgentMgr) OnQuit() {
	// 通知所有room强制存储并退出
	//for _, room := range self.id2Room {
	//	RunOnRoom(room.msgsFromMgr, room, func(input *Room) {
	//		input.SaveAndQuit()
	//	})
	//}
	//
	//for len(self.id2Room) > 0 {
	//	c := <-self.MsgsFromRoom
	//	SafeRunClosure(self, c)
	//}

}

func (self *AgentMgr) Quit() {
	close(self.MsgFromAgent)
}

func (self *AgentMgr) SetNextFrameId() {
	self.frameID++
	if self.frameID < 0 {
		self.frameID = 1
	}
}

func (self *AgentMgr) EnterGame(agent *Agent) {
	oldAgent := self.uid2Agent[agent.Uid]
	if oldAgent != nil {
		//踢掉老链接
		oldAgent.CloseConnect()
		delete(self.uid2Agent, agent.Uid)
	}

	self.uid2Agent[agent.Uid] = agent
}

func (self *AgentMgr) loseConnect(agent *Agent) {
	delete(self.uid2Agent, agent.Uid)
}

func (self *AgentMgr) MatchResponse(res *pb.MatchOverRes) {
	//找出对应的agent
	agent := self.uid2Agent[res.GetUid()]
	if agent == nil {
		self.Log.Error("MatchResponse uid %+v not exist", res.GetUid())
		//找不到agent, 说明已经断开链接
		return
	}

	RunOnAgent(agent.MsgFromAgentMgr, agent, func(agent *Agent) {
		agent.MatchResponse(res)
	})

}
