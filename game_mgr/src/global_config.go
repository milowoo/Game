package game_mgr

import (
	"game_mgr/src/domain"
	"game_mgr/src/internal"
	"github.com/go-ini/ini"
)

type GlobalConfig struct {
	Port        int
	Level       string
	LocalModel  bool
	RedisConfig *domain.RedisConfig
	MongoConfig *domain.MongoConfig
	NatsConfig  *domain.NatsConfig
	NacosConfig *domain.NacosConfig
}

const (
	CFG_NAME = "/Users/wuchuangeng/game/game_mgr/conf/game.ini"
)

func NewGlobalConfig() (*GlobalConfig, error) {
	log := internal.GLog
	ret := &GlobalConfig{}
	cfg, err := ini.Load(CFG_NAME)
	if err != nil {
		return nil, err
	}

	section, err := cfg.GetSection("game_mgr")
	if err != nil {
		log.Error("NewGlobalConfig  GetSection game_mgr  err")
		return nil, err
	}

	key, err := section.GetKey("level")
	if err != nil {
		log.ResetLevel("info")
	} else {
		logLevel := key.String()
		log.ResetLevel(logLevel)
	}

	key, err = section.GetKey("port")
	if err != nil {
		log.Error("get section game_mgr key port  err ")
		return nil, err
	}
	ret.Port = key.MustInt()

	key, err = section.GetKey("localMode")
	if err != nil {
		log.Error("get section game_mgr key localMode  err ")
		return nil, err
	}
	ret.LocalModel = key.MustBool(false)

	nacosConfig, err := LoadNacosConfig()
	if err != nil {
		log.Error("get section nacos   err ")
		return nil, err
	}

	ret.NacosConfig = nacosConfig

	mongoConfig, err := LoadMongoConfig()

	ret.MongoConfig = mongoConfig

	redisConfig, err := LoadRedisConfig()
	if err != nil {
		log.Error("get section redis   err ")
		return nil, err
	}

	ret.RedisConfig = redisConfig

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

func LoadNacosConfig() (*domain.NacosConfig, error) {
	log := internal.GLog
	ret := &domain.NacosConfig{}
	cfg, err := ini.Load(CFG_NAME)
	if err != nil {
		log.Error("load file game.ini err ")
		return nil, err
	}

	section, err := cfg.GetSection("nacos")
	if err != nil {
		log.Error("get redis  err ")
		return nil, err
	}

	key, err := section.GetKey("ip")
	if err != nil {
		log.Error("get room section key Ip  err ")
		return nil, err
	}
	ret.Ip = key.MustString("127.0.0.1")

	key, err = section.GetKey("port")
	if err != nil {
		log.Error("get room section key port  err ")
		return nil, err
	}
	ret.Port = key.MustUint64(8848)

	key, err = section.GetKey("spaceId")
	if err != nil {
		log.Error("get room section key spaceId  err ")
		return nil, err
	}
	ret.SpaceId = key.MustString("xxxx")

	key, err = section.GetKey("gameGroup")
	if err != nil {
		log.Error("get room section key gameGroup  err ")
		return nil, err
	}
	ret.GameGroup = key.MustString("gameGroup")

	key, err = section.GetKey("gameDataId")
	if err != nil {
		log.Error("get room section key allGameIdGroup  err ")
		return nil, err
	}
	ret.GameDataId = key.MustString("gameDataId")

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
