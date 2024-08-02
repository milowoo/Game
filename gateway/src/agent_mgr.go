package gateway

import (
	"gateway/src/internal"
	"gateway/src/pb"
	"time"
)

const (
	RoomMgrFrameInterval = time.Millisecond * 500
)

type AgentMgr struct {
	Server       *Server
	uid2Agent    map[string]*Agent
	MsgFromAgent chan Closure
	frameTimer   *time.Ticker
	frameID      int
}

func NewAgentMgr(server *Server) *AgentMgr {
	return &AgentMgr{
		Server:       server,
		uid2Agent:    make(map[string]*Agent),
		MsgFromAgent: make(chan Closure, 16*1024),
		frameTimer:   time.NewTicker(RoomMgrFrameInterval),
		frameID:      0,
	}
}

func (self *AgentMgr) Run() {
	defer func() {
		p := recover()
		if p != nil {
			internal.GLog.Info("execute panic recovered and going to stop: %v", p)
		}
	}()

	internal.GLog.Info("agent mgr begin ....")

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
	//	self.Log.Info("frame begin ...")
	self.SetNextFrameId()
}

func (self *AgentMgr) OnQuit() {
	//通知所有room强制存储并退出
	for _, agent := range self.uid2Agent {
		RunOnAgent(agent.MsgFromAgentMgr, agent, func(agent *Agent) {
			agent.OnClose()
		})
	}

}

func (self *AgentMgr) Quit() {
	internal.GLog.Info("agent mgr quit ...")
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
	internal.GLog.Info("agent mgr lose connect uid %+v ....", agent.Uid)
	delete(self.uid2Agent, agent.Uid)
}

func (self *AgentMgr) MatchResponse(res *pb.MatchOverRes) {
	//找出对应的agent
	agent := self.uid2Agent[res.GetUid()]
	if agent == nil {
		internal.GLog.Error("MatchResponse uid %+v not exist", res.GetUid())
		//找不到agent, 说明已经断开链接
		return
	}

	RunOnAgent(agent.MsgFromAgentMgr, agent, func(agent *Agent) {
		agent.CallMatchResponse(res)
	})

}

func (self *AgentMgr) GamePushDataNotice(res *pb.GamePushMessage) {
	//找出对应的agent
	head := res.GetHead()
	agent := self.uid2Agent[head.GetUid()]
	if agent == nil {
		internal.GLog.Error("GamePushDataNotice uid %+v not exist", head.GetUid())
		//找不到agent, 说明已经断开链接
		return
	}

	//校验是否是同一款游戏
	if agent.GameId != head.GetGameId() {
		internal.GLog.Error("GamePushDataNotice uid %+v game no equal agent gameId %+v push gameId %+v ",
			head.GetUid(), agent.GameId, head.GameId)
		//找不到agent, 说明已经断开链接
		return
	}

	RunOnAgent(agent.MsgFromAgentMgr, agent, func(agent *Agent) {
		agent.ProcGamePushMessage(res)
	})

}
