package config

import (
	"fmt"
	"gateway/src/domain"
	"gateway/src/internal"
	"github.com/go-ini/ini"
	"time"
)

type GlobalConfig struct {
	Client      *domain.NetworkConfig
	AgentConfig *domain.AgentConfig
	RedisConfig *domain.RedisConfig
	NatsConfig  *domain.NatsConfig
	NacosConfig *domain.NacosConfig

	MaxClient int
}

const (
	CFG_NAME = "/Users/wuchuangeng/game/gateway/conf/game.ini"
)

func NewGlobalConfig() (*GlobalConfig, error) {
	ret := &GlobalConfig{}
	cfg, err := ini.Load(CFG_NAME)
	if err != nil {
		internal.GLog.Error("load file game.ini err ")
		return nil, err
	}

	section, err := cfg.GetSection("gateway")
	if err != nil {
		internal.GLog.Error("get room section  err ")
		return nil, err
	}

	key, err := section.GetKey("level")
	if err != nil {
		internal.GLog.ResetLevel("info")
	} else {
		logLevel := key.String()
		internal.GLog.ResetLevel(logLevel)
	}

	key, err = section.GetKey("maxClient")
	if err != nil {
		internal.GLog.Error("get  section key maxClient  err ")
		return nil, err
	}
	ret.MaxClient = key.MustInt()

	ret.RedisConfig, err = NewRedisConfig()
	if err != nil {
		internal.GLog.Error("get redis section key  err ")
		return nil, err
	}

	ret.AgentConfig, err = NewAgentConfig()
	if err != nil {
		internal.GLog.Error("get agent section key  err ")
		return nil, err
	}

	ret.NatsConfig, err = NewNatsConfig()
	if err != nil {
		internal.GLog.Error("get nats section key  err ")
		return nil, err
	}

	ret.NacosConfig, err = LoadNacosConfig()
	if err != nil {
		internal.GLog.Error("get nacos section key  err ")
		return nil, err
	}

	ret.Client, err = LoadClientConfig()
	if err != nil {
		internal.GLog.Error("get client section key  err ")
		return nil, err
	}

	internal.GLog.Info("NewGlobalConfig  success")

	return ret, nil
}

func LoadClientConfig() (*domain.NetworkConfig, error) {
	ret := &domain.NetworkConfig{}
	cfg, err := ini.Load(CFG_NAME)
	if err != nil {
		internal.GLog.Error("load file game.ini err ")
		return nil, err
	}

	section, err := cfg.GetSection("client")
	if err != nil {
		internal.GLog.Error("get client  err ")
		return nil, err
	}

	key, err := section.GetKey("listenIp")
	if err != nil {
		internal.GLog.Error("get client section key listenIp  err ")
		return nil, err
	}
	ret.ListenIp = key.MustString("127.0.0.1")

	key, err = section.GetKey("port")
	if err != nil {
		internal.GLog.Error("get client section key masterName  err ")
		return nil, err
	}
	ret.ListenPort = key.MustInt(36001)

	key, err = section.GetKey("timeout")
	if err != nil {
		internal.GLog.Error("get client section key timeout  err ")
		return nil, err
	}
	ret.Timeout = key.MustInt(10)

	key, err = section.GetKey("hmacKey")
	if err != nil {
		internal.GLog.Error("get client section key hmacKey  err ")
		return nil, err
	}
	ret.HMACKey = key.MustString("xxx")

	return ret, nil
}

func NewAgentConfig() (*domain.AgentConfig, error) {
	ret := &domain.AgentConfig{}
	cfg, err := ini.Load(CFG_NAME)
	if err != nil {
		internal.GLog.Error("load file game.ini err ")
		return nil, err
	}

	section, err := cfg.GetSection("agent")
	if err != nil {
		internal.GLog.Error("get agent  err ")
		return nil, err
	}

	key, err := section.GetKey("enableLogSend")
	if err != nil {
		internal.GLog.Error("get agent key enableLogSend err ")
		return nil, err
	}
	ret.EnableLogSend = key.MustBool(false)

	key, err = section.GetKey("enableLogRecv")
	if err != nil {
		internal.GLog.Error("get room agent key enableLogRecv err ")
		return nil, err
	}
	ret.EnableLogRecv = key.MustBool(false)

	key, err = section.GetKey("enableCheckPing")
	if err != nil {
		internal.GLog.Error("get agent key enableCheckPing err ")
		return nil, err
	}
	ret.EnableCheckPing = key.MustBool(true)

	key, err = section.GetKey("enableCheckLoginParams")
	if err != nil {
		internal.GLog.Error("get agent key enableCheckLoginParams err ")
		return nil, err
	}
	ret.EnableCheckLoginParams = key.MustBool(true)

	key, err = section.GetKey("enableCachedMsg")
	if err != nil {
		internal.GLog.Error("get agent key enableCachedMsg err ")
		return nil, err
	}
	ret.EnableCachedMsg = key.MustBool(true)

	key, err = section.GetKey("enableCompressMsg")
	if err != nil {
		internal.GLog.Error("get agent key enableCompressMsg err ")
		return nil, err
	}
	ret.EnableCompressMsg = key.MustBool(false)

	ret.CachedMsgMaxCount, err = LoadPositiveIntFromIniSection(section, "cachedMsgMaxCount")
	if err != nil {
		internal.GLog.Error("get agent key cachedMsgMaxCount err ")
		return nil, err
	}

	ret.CompressMsgSizeThreshold, err = LoadPositiveIntFromIniSection(section, "compressMsgSizeThreshold")
	if err != nil {
		internal.GLog.Error("get agent key compressMsgSizeThreshold err ")
		return nil, err
	}

	ret.TimestampExpireDuration, err = LoadSecondFromIniSection(section, "timestampExpireSeconds")
	if err != nil {
		internal.GLog.Error("get room agent key timestampExpireSeconds err ")
		return nil, err
	}

	return ret, nil
}

func NewRedisConfig() (*domain.RedisConfig, error) {
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
		internal.GLog.Error("get redis section key maxClient  err ")
		return nil, err
	}
	ret.Address = key.MustString("127.0.0.1:6377")

	key, err = section.GetKey("masterName")
	if err != nil {
		internal.GLog.Error("get redis section key masterName  err ")
		return nil, err
	}
	ret.MasterName = key.MustString("")

	key, err = section.GetKey("password")
	if err != nil {
		internal.GLog.Error("get redis section key password  err ")
		return nil, err
	}
	ret.Password = key.MustString("")

	return ret, nil
}

func NewNatsConfig() (*domain.NatsConfig, error) {
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

func LoadPositiveIntFromIniSection(section *ini.Section, keyName string) (int, error) {
	key, err := section.GetKey(keyName)
	if err != nil {
		return -1, err
	}

	v := key.MustInt(-1)
	if v <= 0 {
		return -1, fmt.Errorf("unexpected %s.%s=%d", section.Name(), keyName, v)
	}

	return v, nil
}

func LoadSecondFromIniSection(section *ini.Section, keyName string) (time.Duration, error) {
	key, err := section.GetKey(keyName)
	if err != nil {
		return -1, err
	}

	f := key.MustFloat64(-1)
	if f <= 0 {
		return -1, fmt.Errorf("unexpected %s.%s=%f", section.Name(), keyName, f)
	}

	return time.Duration(float64(time.Second) * f), nil
}
