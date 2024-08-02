package config

import (
	"game_frame/src/domain"
	"game_frame/src/internal"
	"github.com/go-ini/ini"
)

type GlobalConfig struct {
	GameId      string
	Port        uint64
	GroupName   string
	Level       string
	RedisConfig *domain.RedisConfig
	NatsConfig  *domain.NatsConfig
	NacosConfig *domain.NacosConfig
	MongoConfig *domain.MongoConfig
}

const (
	CFG_NAME = "/Users/wuchuangeng/game/game_frame/conf/game.ini"
)

func NewGlobalConfig() (*GlobalConfig, error) {
	ret := &GlobalConfig{}
	cfg, err := ini.Load(CFG_NAME)
	if err != nil {
		return nil, err
	}

	section, err := cfg.GetSection("game")
	if err != nil {
		internal.GLog.Error("NewGlobalConfig  GetSection game  err")
		return nil, err
	}

	key, err := section.GetKey("gameId")
	if err != nil {
		internal.GLog.Error("NewGlobalConfig  GetSection game  gameId err")
		return nil, err
	}

	ret.GameId = key.MustString("xx")

	key, err = section.GetKey("port")
	if err != nil {
		internal.GLog.Error("NewGlobalConfig  GetSection game  port err")
		return nil, err
	}

	ret.Port = key.MustUint64(7723)

	key, err = section.GetKey("gameGroup")
	if err != nil {
		internal.GLog.Error("NewGlobalConfig  GetSection game  gameId err")
		return nil, err
	}

	ret.GroupName = key.MustString("A")

	key, err = section.GetKey("level")
	if err != nil {
		internal.GLog.ResetLevel("info")
	} else {
		logLevel := key.String()
		internal.GLog.ResetLevel(logLevel)
	}

	redisConfig, err := LoadRedisConfig()
	if err != nil {
		internal.GLog.Error("get section redis  err ")
		return nil, err
	}

	ret.RedisConfig = redisConfig

	ret.MongoConfig, err = LoadMongoConfig()
	if err != nil {
		internal.GLog.Error("get section mongo err ")
		return nil, err
	}

	ret.NatsConfig, err = LoadNatsConfig()
	if err != nil {
		internal.GLog.Error("get nats section key  err ")
		return nil, err
	}

	ret.NacosConfig, err = LoadNacosConfig()
	if err != nil {
		internal.GLog.Error("get nacos section key  err ")
		return nil, err
	}

	return ret, nil
}

func LoadNatsConfig() (*domain.NatsConfig, error) {
	ret := &domain.NatsConfig{}
	cfg, err := ini.Load(CFG_NAME)
	if err != nil {
		internal.GLog.Error("load file game.ini err ")
		return nil, err
	}

	section, err := cfg.GetSection("nats")
	if err != nil {
		internal.GLog.Error("get redis  err ")
		return nil, err
	}

	key, err := section.GetKey("address")
	if err != nil {
		internal.GLog.Error("get room section key maxClient  err ")
		return nil, err
	}
	ret.Address = key.MustString("127.0.0.1:6377")
	return ret, nil
}

func LoadRedisConfig() (*domain.RedisConfig, error) {
	ret := &domain.RedisConfig{}
	cfg, err := ini.Load(CFG_NAME)
	if err != nil {
		internal.GLog.Error("load file game.ini err ")
		return nil, err
	}

	section, err := cfg.GetSection("redis")
	if err != nil {
		internal.GLog.Error("get redis  err ")
		return nil, err
	}

	key, err := section.GetKey("address")
	if err != nil {
		internal.GLog.Error("get room section key maxClient  err ")
		return nil, err
	}
	ret.Address = key.MustString("127.0.0.1:6377")

	key, err = section.GetKey("masterName")
	if err != nil {
		internal.GLog.Error("get room section key maxClient  err ")
		return nil, err
	}
	ret.MasterName = key.MustString("")

	key, err = section.GetKey("password")
	if err != nil {
		internal.GLog.Error("get room section key password  err ")
		return nil, err
	}
	ret.Password = key.MustString("")

	return ret, nil
}

func LoadMongoConfig() (*domain.MongoConfig, error) {
	ret := &domain.MongoConfig{}
	cfg, err := ini.Load(CFG_NAME)
	if err != nil {
		internal.GLog.Error("load file game.ini err ")
		return nil, err
	}

	section, err := cfg.GetSection("mongo")
	if err != nil {
		internal.GLog.Error("get redis  err ")
		return nil, err
	}

	key, err := section.GetKey("address")
	if err != nil {
		internal.GLog.Error("get room section key maxClient  err ")
		return nil, err
	}
	ret.Address = key.MustString("127.0.0.1:10009")

	key, err = section.GetKey("dbname")
	if err != nil {
		internal.GLog.Error("get room section key dbname  err ")
		return nil, err
	}
	ret.Name = key.MustString("")

	return ret, nil
}

func LoadNacosConfig() (*domain.NacosConfig, error) {
	ret := &domain.NacosConfig{}
	cfg, err := ini.Load(CFG_NAME)
	if err != nil {
		internal.GLog.Error("load file game.ini err ")
		return nil, err
	}

	section, err := cfg.GetSection("nacos")
	if err != nil {
		internal.GLog.Error("get nacos  err ")
		return nil, err
	}

	key, err := section.GetKey("ip")
	if err != nil {
		internal.GLog.Error("get nacos section key Ip  err ")
		return nil, err
	}
	ret.Ip = key.MustString("127.0.0.1")

	key, err = section.GetKey("port")
	if err != nil {
		internal.GLog.Error("get nacos section key port  err ")
		return nil, err
	}
	ret.Port = key.MustUint64(8848)

	key, err = section.GetKey("spaceId")
	if err != nil {
		internal.GLog.Error("get nacos section key spaceId  err ")
		return nil, err
	}
	ret.SpaceId = key.MustString("xxxx")

	key, err = section.GetKey("gameGroup")
	if err != nil {
		internal.GLog.Error("get nacos section key gameGroup  err ")
		return nil, err
	}
	ret.GameGroup = key.MustString("gameGroup")

	key, err = section.GetKey("gameDataId")
	if err != nil {
		internal.GLog.Error("get nacos section key gameDataId  err ")
		return nil, err
	}
	ret.GameDataId = key.MustString("gameDataId")

	return ret, nil
}
