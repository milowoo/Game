package internal

import (
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"match/src/log"
	"match/src/mq"
	"match/src/redis"
)

var GLog *log.Logger
var RedisDao *redis.RedisDao
var NatsPool *mq.NatsPool
var ConfigClient config_client.IConfigClient
var NameClient naming_client.INamingClient
