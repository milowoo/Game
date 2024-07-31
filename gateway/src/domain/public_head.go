package domain

type PublicHead struct {
	GameId    string `json:"gameId,omitempty"`
	RoomId    string `json:"roomId,omitempty"`
	Uid       string `json:"uid,omitempty"`
	Pid       string `json:"pid,omitempty"`
	Sn        int64  `json:"sn,omitempty"`
	Timestamp int64  `json:"timestamp,omitempty"`
	PbName    string `json:"pbName,omitempty"`
	GatewayIp string `json:"gatewayIp,omitempty"`
	Data      string `json:"data,omitempty"`
}
