package mongo

import (
	"context"
	"game_mgr/src/domain"
	"game_mgr/src/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
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

// 获取游戏信息
func (dao *MongoDAO) GetGame(gameId string) *domain.GameInfo {
	c := dao.database.Collection("game_info")
	var m domain.GameInfo
	err := c.FindOne(context.Background(), bson.M{"gameid": gameId}).Decode(&m)
	if err != nil {
		dao.log.Info("GetGame err %v gameId %v ", err, gameId)
		return nil
	}
	return &m
}

func (dao *MongoDAO) InsertGame(game *domain.GameInfo) error {
	game.UTime = time.Now()
	game.CTime = time.Now()

	c := dao.database.Collection("game_info")
	_, err := c.InsertOne(context.Background(), game)
	if err != nil {
		dao.log.Error("InsertGame err %v gameId %v", game.GameId)
		return err
	}

	return nil

}

func (dao *MongoDAO) SavaGame(gameId string, game *domain.GameInfo) error {
	game.UTime = time.Now()

	query := bson.M{"gameid": gameId}
	updateFilter := bson.M{"$set": bson.M{
		"name":      game.Name,
		"status":    game.Status,
		"groupname": game.GroupName,
		"gametime":  game.GameTime,
		"matchtime": game.MatchTime,
		"operator":  game.Operator,
		"utime":     game.UTime}}
	c := dao.database.Collection("game_info")
	_, err := c.UpdateOne(context.Background(), query, updateFilter)
	if err != nil {
		dao.log.Info("SavaGame err %v, gameId %v", err, gameId)
		return err
	}

	return nil
}
