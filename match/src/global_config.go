package match

import (
	"fmt"
	"github.com/go-ini/ini"
	"match/src/log"
)

type GlobalConfig struct {
	Level       string
	RedisConfig *RedisConfig
	NatsConfig  *NatsConfig
	NacosConfig *NacosConfig
}

type RedisConfig struct {
	Address    string
	MasterName string
	Password   string
}

type NatsConfig struct {
	Address string
}

type NacosConfig struct {
	Ip         string
	Port       uint64
	SpaceId    string
	GameGroup  string
	GameDataId string
}

const (
	CFG_DIR = "../conf"
)

func NewGlobalConfig(log *log.Logger) (*GlobalConfig, error) {
	ret := &GlobalConfig{}
	cfg, err := ini.Load(fmt.Sprintf("%s/game.ini", CFG_DIR))
	if err != nil {
		return nil, err
	}

	section, err := cfg.GetSection("match")
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

	redisConfig, err := LoadRedisConfig(log)
	if err != nil {
		log.Error("get section redis   err ")
		return nil, err
	}

	ret.RedisConfig = redisConfig

	nacosConfig, err := LoadNacosConfig(log)
	if err != nil {
		log.Error("get section nacos err ")
		return nil, err
	}

	ret.NacosConfig = nacosConfig

	ret.NatsConfig, err = LoadNatsConfig(log)
	if err != nil {
		log.Error("get nats section key  err ")
		return nil, err
	}

	return ret, nil
}

func LoadNacosConfig(log *log.Logger) (*NacosConfig, error) {
	ret := &NacosConfig{}
	cfg, err := ini.Load(fmt.Sprintf("%s/game.ini", "../conf"))
	if err != nil {
		log.Error("load file game.ini err ")
		return nil, err
	}

	section, err := cfg.GetSection("nacos")
	if err != nil {
		log.Error("get nacos  err ")
		return nil, err
	}

	key, err := section.GetKey("Ip")
	if err != nil {
		log.Error("get nacos section key Ip  err ")
		return nil, err
	}
	ret.Ip = key.MustString("127.0.0.1")

	key, err = section.GetKey("port")
	if err != nil {
		log.Error("get nacos section key port  err ")
		return nil, err
	}
	ret.Port = key.MustUint64(8848)

	key, err = section.GetKey("spaceId")
	if err != nil {
		log.Error("get nacos section key spaceId  err ")
		return nil, err
	}
	ret.SpaceId = key.MustString("xxxx")

	key, err = section.GetKey("gameGroup")
	if err != nil {
		log.Error("get nacos section key gameGroup  err ")
		return nil, err
	}
	ret.GameGroup = key.MustString("gameGroup")

	key, err = section.GetKey("gameDataId")
	if err != nil {
		log.Error("get nacos section key gameDataId  err ")
		return nil, err
	}
	ret.GameDataId = key.MustString("gameDataId")

	return ret, nil
}

func LoadNatsConfig(log *log.Logger) (*NatsConfig, error) {
	ret := &NatsConfig{}
	cfg, err := ini.Load(fmt.Sprintf("%s/game.ini", "../conf"))
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

func LoadRedisConfig(log *log.Logger) (*RedisConfig, error) {
	ret := &RedisConfig{}
	cfg, err := ini.Load(fmt.Sprintf("%s/game.ini", "../conf"))
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
