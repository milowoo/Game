package domain

import (
	"time"
)

/*
**	player 玩家记录
 */
type Player struct {
	Uid   string    `bson:"uid,omitempty"`
	Pid   string    `bson:"pid" json:"pid"`
	Utime time.Time `bson:"utime" json:"utime"`
	Ctime time.Time `bson:"ctime" json:"ctime"`
}

func NewPlayer(pid string, uid string) *Player {
	data := &Player{Pid: pid,
		Uid: uid}

	return data
}
