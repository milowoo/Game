package internal

import (
	"game_mgr/src/log"
	"game_mgr/src/mongo"
	"game_mgr/src/mq"
	"game_mgr/src/redis"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
)

var GLog *log.Logger
var RedisDao *redis.RedisDao
var Mongo *mongo.MongoDAO
var NatsPool *mq.NatsPool
var IClient config_client.IConfigClient
