package mongo

import (
	"game_mgr/src/domain"
	"game_mgr/src/log"
	"time"

	"gopkg.in/mgo.v2/bson"
)

type MongoDAO struct {
	BaseDAO
	log *log.Logger
}

func NewMongoDao(dataSource *DataSource, log *log.Logger) *MongoDAO {
	dao := &MongoDAO{BaseDAO: *NewBaseDAO(dataSource),
		log: log}
	return dao
}

// 获取游戏信息
func (dao *MongoDAO) GetGame(gameId string) (*domain.GameInfo, error) {
	session := dao.dataSource.GetSession()
	defer session.Close()

	c := session.DB(dao.dataSource.database).C("game")
	m := domain.GameInfo{}
	err := c.Find(&bson.M{"_id": gameId}).One(&m)
	if err != nil {
		dao.log.Debug("GetGame err %v uid %v ", err, gameId)
		return &m, err
	}
	return &m, nil
}

func (dao *MongoDAO) InsertGame(game *domain.GameInfo) error {
	game.UTime = time.Now()
	game.CTime = time.Now()

	session := dao.dataSource.GetSession()
	defer session.Close()
	c := session.DB(dao.dataSource.database).C("game")
	err := c.Insert(game)
	if err != nil {
		dao.log.Error("InsertGame err %v gameId %v", game.GameId)
		return err
	}

	return nil

}

func (dao *MongoDAO) SavaGame(gameId string, game interface{}) error {
	session := dao.dataSource.GetSession()
	defer session.Close()
	c := session.DB(dao.dataSource.database).C("game")

	query := bson.M{"_id": gameId}

	_, err := c.Upsert(query, game)
	if err != nil {
		dao.log.Debug("UpdatePlayer err %v, uid %v", err, gameId)
		return err
	}

	return nil
}

func (dao *MongoDAO) DeleteGame(gameId string) error {
	session := dao.dataSource.GetSession()
	defer session.Close()
	c := session.DB(dao.dataSource.database).C("game")

	err := c.Remove(bson.M{"_id": gameId})

	if err != nil {
		dao.log.Debug("DeleteGame err %v gameId %v", err, gameId)
		return err
	}

	return nil
}

func (dao *MongoDAO) GetGameByTime(updateTime int64) ([]domain.GameInfo, error) {
	session := dao.dataSource.GetSession()
	defer session.Close()
	m := make([]domain.GameInfo, 0)

	c := session.DB(dao.dataSource.database).C("game")
	var game domain.GameInfo
	iter := c.Find(bson.M{"uTime": bson.M{"$gt": updateTime}}).Iter()
	for iter.Next(&game) {
		m = append(m, game)
	}

	return m, nil
}
