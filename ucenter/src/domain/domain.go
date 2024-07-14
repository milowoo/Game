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
	Url   string        `bson:"url" json:"url"`
	Sign  string        `bson:"sign" json:"sign"`
	Sex   int8          `bson:"sex" json:"sex"`
	Utime time.Time     `bson:"utime" json:"utime"`
}

func NewPlayer(uid string) *Player {
	data := &Player{Uid: bson.ObjectId(uid)}

	return data
}
