package internal

import (
	"gateway/src/log"
	"gateway/src/mq"
	"gateway/src/redis"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
)

var GLog *log.Logger
var RedisDao *redis.RedisDao
var ConfigClient config_client.IConfigClient
var NatsPool *mq.NatsPool
