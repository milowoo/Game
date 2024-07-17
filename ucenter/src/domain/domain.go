package domain

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

/*
**	player 玩家记录
 */
type Player struct {
	Uid   bson.ObjectId `bson:"_id,omitempty"`
	Pid   string        `bson:"sign" json:"sign"`
	Utime time.Time     `bson:"utime" json:"utime"`
}

func NewPlayer(uid string) *Player {
	data := &Player{Uid: bson.ObjectId(uid)}

	return data
}
