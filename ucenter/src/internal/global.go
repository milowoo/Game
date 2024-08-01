package internal

import (
	"ucenter/src/log"
	"ucenter/src/mongo"
	"ucenter/src/mq"
	"ucenter/src/redis"
)

var GLog *log.Logger
var RedisDao *redis.RedisDao
var Mongo *mongo.MongoDAO
var NatsPool *mq.NatsPool
