package domain

type GamePlayer struct {
	Uid             string
	Pid             string
	RoomId          string
	HallId          string
	GatewayIp       string
	LoadingProgress int32
	OffLineFrameId  int
	TotalUseTime    int
	Score           int
	IsAi            bool
	IsNewPlayer     bool
}
