package ucenter

import (
	"fmt"
	"github.com/go-ini/ini"
	"ucenter/src/log"
)

type GlobalConfig struct {
	Level       string
	RedisConfig *RedisConfig
	NatsConfig  *NatsConfig
	MongoConfig *MongoConfig
}

type RedisConfig struct {
	Address    string
	MasterName string
	Password   string
}

type NatsConfig struct {
	Address string
}

type MongoConfig struct {
	Address string
	Name    string
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

	redisConfig, err := LoadRedisConfig(log)
	if err != nil {
		log.Error("get section redis  err ")
		return nil, err
	}

	ret.RedisConfig = redisConfig

	ret.MongoConfig, err = LoadMongoConfig(log)
	if err != nil {
		log.Error("get section mongo err ")
		return nil, err
	}

	ret.NatsConfig, err = LoadNatsConfig(log)
	if err != nil {
		log.Error("get nats section key  err ")
		return nil, err
	}

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

func LoadMongoConfig(log *log.Logger) (*MongoConfig, error) {
	ret := &MongoConfig{}
	cfg, err := ini.Load(fmt.Sprintf("%s/game.ini", "../conf"))
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
