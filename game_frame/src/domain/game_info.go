package domain

import "time"

type GameInfo struct {
	GameId    string    `json:"gameId,omitempty"`
	Type      int32     `json:"type,omitempty"`
	Name      string    `json:"name,omitempty"`
	GroupName string    `json:"groupName,omitempty"`
	Status    int32     `json:"status,omitempty"`
	GameTime  int32     `json:"gameTime,omitempty"`
	MatchTime int32     `json:"matchTime,omitempty"`
	Operator  string    `json:"operator,omitempty"`
	UTime     time.Time `bson:"utime" json:"utime"`
	CTime     time.Time `bson:"ctime" json:"ctime"`
}
