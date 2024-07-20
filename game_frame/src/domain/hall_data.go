package domain

type HallData struct {
	GameId string `json:"gameId,omitempty"`
	Uid    string `json:"uid,omitempty"`
	RoomId string `json:"roomId,omitempty"`
	UTime  int64  `json:"uTime,omitempty"`
}
