package gateway

import (
	"fmt"
	"gateway/src/log"
	"github.com/go-ini/ini"
	"time"
)

type GlobalConfig struct {
	Client      *NetworkAsServer
	AgentConfig *AgentConfig
	RedisConfig *RedisConfig
	NatsConfig  *NatsConfig
	NacosConfig *NacosConfig

	MaxClient int
}

type NetworkAsServer struct {
	ListenIp   string
	ListenPort int
	Timeout    time.Duration
	HMACKey    string
}

type AgentConfig struct {
	EnableLogSend          bool
	EnableLogRecv          bool
	EnableCheckPing        bool
	EnableCheckLoginParams bool

	EnableCachedMsg   bool
	EnableCompressMsg bool

	CachedMsgMaxCount        int
	CompressMsgSizeThreshold int

	TimestampExpireDuration time.Duration
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

func NewGlobalConfig(log *log.Logger) (*GlobalConfig, error) {
	ret := &GlobalConfig{}
	cfg, err := ini.Load(fmt.Sprintf("%s/game.ini", "../conf"))
	if err != nil {
		log.Error("load file game.ini err ")
		return nil, err
	}

	section, err := cfg.GetSection("gateway")
	if err != nil {
		log.Error("get room section  err ")
		return nil, err
	}

	key, err := section.GetKey("level")
	if err != nil {
		log.ResetLevel("info")
	} else {
		logLevel := key.String()
		log.ResetLevel(logLevel)
	}

	key, err = section.GetKey("maxClient")
	if err != nil {
		log.Error("get  section key maxClient  err ")
		return nil, err
	}
	ret.MaxClient = key.MustInt()

	ret.RedisConfig, err = NewRedisConfig(log)
	if err != nil {
		log.Error("get redis section key  err ")
		return nil, err
	}

	ret.AgentConfig, err = NewAgentConfig(log)
	if err != nil {
		log.Error("get agent section key  err ")
		return nil, err
	}

	ret.NatsConfig, err = NewNatsConfig(log)
	if err != nil {
		log.Error("get nats section key  err ")
		return nil, err
	}

	ret.NacosConfig, err = LoadNacosConfig(log)
	if err != nil {
		log.Error("get nacos section key  err ")
		return nil, err
	}

	log.Info("NewGlobalConfig  success")

	return ret, nil
}

func NewAgentConfig(log *log.Logger) (*AgentConfig, error) {
	ret := &AgentConfig{}
	cfg, err := ini.Load(fmt.Sprintf("%s/game.ini", "../conf"))
	if err != nil {
		log.Error("load file game.ini err ")
		return nil, err
	}

	section, err := cfg.GetSection("agent")
	if err != nil {
		log.Error("get agent  err ")
		return nil, err
	}

	key, err := section.GetKey("enableLogSend")
	if err != nil {
		log.Error("get agent key enableLogSend err ")
		return nil, err
	}
	ret.EnableLogSend = key.MustBool(false)

	key, err = section.GetKey("enableLogRecv")
	if err != nil {
		log.Error("get room agent key enableLogRecv err ")
		return nil, err
	}
	ret.EnableLogRecv = key.MustBool(false)

	key, err = section.GetKey("enableCheckPing")
	if err != nil {
		log.Error("get agent key enableCheckPing err ")
		return nil, err
	}
	ret.EnableCheckPing = key.MustBool(true)

	key, err = section.GetKey("enableCheckLoginParams")
	if err != nil {
		log.Error("get agent key enableCheckLoginParams err ")
		return nil, err
	}
	ret.EnableCheckLoginParams = key.MustBool(true)

	key, err = section.GetKey("enableCachedMsg")
	if err != nil {
		log.Error("get agent key enableCachedMsg err ")
		return nil, err
	}
	ret.EnableCachedMsg = key.MustBool(true)

	key, err = section.GetKey("enableCompressMsg")
	if err != nil {
		log.Error("get agent key enableCompressMsg err ")
		return nil, err
	}
	ret.EnableCompressMsg = key.MustBool(false)

	ret.CachedMsgMaxCount, err = LoadPositiveIntFromIniSection(section, "cachedMsgMaxCount")
	if err != nil {
		log.Error("get agent key cachedMsgMaxCount err ")
		return nil, err
	}

	ret.CompressMsgSizeThreshold, err = LoadPositiveIntFromIniSection(section, "compressMsgSizeThreshold")
	if err != nil {
		log.Error("get agent key compressMsgSizeThreshold err ")
		return nil, err
	}

	ret.TimestampExpireDuration, err = LoadSecondFromIniSection(section, "timestampExpireSeconds")
	if err != nil {
		log.Error("get room agent key timestampExpireSeconds err ")
		return nil, err
	}

	return ret, nil
}

func NewRedisConfig(log *log.Logger) (*RedisConfig, error) {
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
		log.Error("get redis section key maxClient  err ")
		return nil, err
	}
	ret.Address = key.MustString("127.0.0.1:6377")

	key, err = section.GetKey("masterName")
	if err != nil {
		log.Error("get redis section key masterName  err ")
		return nil, err
	}
	ret.MasterName = key.MustString("")

	key, err = section.GetKey("password")
	if err != nil {
		log.Error("get redis section key password  err ")
		return nil, err
	}
	ret.Password = key.MustString("")

	return ret, nil
}

func NewNatsConfig(log *log.Logger) (*NatsConfig, error) {
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

func NewNetworkAsServerFromIniFile(cfg *ini.File, sectionName string) (*NetworkAsServer, error) {
	if cfg == nil {
		return nil, fmt.Errorf("cfg is nil")
	}

	network := &NetworkAsServer{
		ListenIp:   "0.0.0.0",
		ListenPort: -1,
		Timeout:    0,
	}

	section, err := cfg.GetSection(sectionName)
	if err != nil {
		return nil, err
	}

	key, err := section.GetKey("listenIp")
	if err != nil {
		return nil, err
	}
	network.ListenIp = key.String()

	key, err = section.GetKey("listenPort")
	if err != nil {
		return nil, err
	}
	network.ListenPort = key.MustInt(-1)
	if network.ListenPort <= 0 {
		return nil, fmt.Errorf("unexpected listenPort=%d", network.ListenPort)
	}

	key, err = section.GetKey("timeout")
	if err != nil {
		return nil, err
	}

	seconds := key.MustFloat64(-1.0)
	if seconds <= 0 {
		return nil, fmt.Errorf("unexpected timeout=%.02f", seconds)
	}
	network.Timeout = time.Duration(seconds * float64(time.Second))

	// optional
	key, err = section.GetKey("hmacKey")
	if err == nil {
		network.HMACKey = key.MustString("")
	}

	return network, nil
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
