package internal

import (
	"game_frame/src/log"
	"game_frame/src/mq"
	"game_frame/src/redis"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
)

var GLog *log.Logger
var RedisDao *redis.RedisDao
var NatsPool *mq.NatsPool
var ConfigClient config_client.IConfigClient
var NameClient naming_client.INamingClient
