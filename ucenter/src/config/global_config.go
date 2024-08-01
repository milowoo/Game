package config

import (
	"github.com/go-ini/ini"
	"ucenter/src/domain"
	"ucenter/src/internal"
)

type GlobalConfig struct {
	Level       string
	RedisConfig *domain.RedisConfig
	NatsConfig  *domain.NatsConfig
	MongoConfig *domain.MongoConfig
}

const (
	CFG_NAME = "/Users/wuchuangeng/game/ucenter/conf/game.ini"
)

func NewGlobalConfig() (*GlobalConfig, error) {
	log := internal.GLog
	ret := &GlobalConfig{}
	cfg, err := ini.Load(CFG_NAME)
	if err != nil {
		return nil, err
	}

	section, err := cfg.GetSection("ucenter")
	if err != nil {
		log.Error("NewGlobalConfig  GetSection ucenter  err")
		return nil, err
	}

	key, err := section.GetKey("level")
	if err != nil {
		log.ResetLevel("info")
	} else {
		logLevel := key.String()
		log.ResetLevel(logLevel)
	}

	redisConfig, err := LoadRedisConfig()
	if err != nil {
		log.Error("get section redis  err ")
		return nil, err
	}

	ret.RedisConfig = redisConfig

	ret.MongoConfig, err = LoadMongoConfig()
	if err != nil {
		log.Error("get section mongo err ")
		return nil, err
	}

	ret.NatsConfig, err = LoadNatsConfig()
	if err != nil {
		log.Error("get nats section key  err ")
		return nil, err
	}

	return ret, nil
}

func LoadNatsConfig() (*domain.NatsConfig, error) {
	log := internal.GLog
	ret := &domain.NatsConfig{}
	cfg, err := ini.Load(CFG_NAME)
	if err != nil {
		log.Error("load file game.ini err ")
		return nil, err
	}

	section, err := cfg.GetSection("nats")
	if err != nil {
		log.Error("get redis  err ")
		return nil, err
	}

	key, err := section.GetKey("address")
	if err != nil {
		log.Error("get room section key maxClient  err ")
		return nil, err
	}
	ret.Address = key.MustString("127.0.0.1:6377")
	return ret, nil
}

func LoadRedisConfig() (*domain.RedisConfig, error) {
	log := internal.GLog
	ret := &domain.RedisConfig{}
	cfg, err := ini.Load(CFG_NAME)
	if err != nil {
		log.Error("load file game.ini err ")
		return nil, err
	}

	section, err := cfg.GetSection("redis")
	if err != nil {
		log.Error("get redis  err ")
		return nil, err
	}

	key, err := section.GetKey("address")
	if err != nil {
		log.Error("get room section key maxClient  err ")
		return nil, err
	}
	ret.Address = key.MustString("127.0.0.1:6377")

	key, err = section.GetKey("masterName")
	if err != nil {
		log.Error("get room section key maxClient  err ")
		return nil, err
	}
	ret.MasterName = key.MustString("")

	key, err = section.GetKey("password")
	if err != nil {
		log.Error("get room section key password  err ")
		return nil, err
	}
	ret.Password = key.MustString("")

	return ret, nil
}

func LoadMongoConfig() (*domain.MongoConfig, error) {
	log := internal.GLog
	ret := &domain.MongoConfig{}
	cfg, err := ini.Load(CFG_NAME)
	if err != nil {
		log.Error("load file game.ini err ")
		return nil, err
	}

	section, err := cfg.GetSection("mongo")
	if err != nil {
		log.Error("get redis  err ")
		return nil, err
	}

	key, err := section.GetKey("address")
	if err != nil {
		log.Error("get room section key maxClient  err ")
		return nil, err
	}
	ret.Address = key.MustString("127.0.0.1:10009")

	key, err = section.GetKey("dbname")
	if err != nil {
		log.Error("get room section key dbname  err ")
		return nil, err
	}
	ret.Name = key.MustString("")

	return ret, nil
}
