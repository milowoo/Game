package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
	"ucenter/src/domain"
	"ucenter/src/log"
)

type MongoDAO struct {
	database *mongo.Database
	log      *log.Logger
}

func NewMongoDao(client *Client, log *log.Logger) *MongoDAO {
	dao := &MongoDAO{database: client.database,
		log: log}
	return dao
}

// 获取设备信息
func (dao *MongoDAO) GetPlayer(pid string) *domain.Player {
	c := dao.database.Collection("player")
	var m domain.Player
	err := c.FindOne(context.Background(), bson.M{"pid": pid}).Decode(&m)
	if err != nil {
		dao.log.Debug("GetPlayer err %+v uid %+v ", err, pid)
		return nil
	}
	return &m
}

func (dao *MongoDAO) InsertPlayer(player *domain.Player) error {
	player.Utime = time.Now()
	player.Ctime = time.Now()

	c := dao.database.Collection("player")
	_, err := c.InsertOne(context.Background(), player)
	if err != nil {
		dao.log.Error("InsertGame err %+v pid %+v", err, player.Pid)
		return err
	}

	return nil

}
