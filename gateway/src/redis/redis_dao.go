package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"strings"
	"time"
)

type RedisDao struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedis(addr string, masterName string, password string) *RedisDao {
	addList := make([]string, 0)
	addSec := strings.Split(addr, ",")
	for _, v := range addSec {
		addList = append(addList, v)
	}

	rdb := redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName:    masterName,
		SentinelAddrs: addList,
		Password:      password,
		DB:            10,
	})

	return &RedisDao{
		client: rdb,
		ctx:    context.Background(),
	}
}

func (s *RedisDao) EXISTS(key string) (bool, error) {
	result, err := s.client.Exists(s.ctx, key).Result()
	return result > 0, err
}

func (s *RedisDao) Get(key string) (string, error) {
	return s.client.Get(s.ctx, key).Result()
}

func (s *RedisDao) Set(key string, value interface{}, expiration time.Duration) error {
	return s.client.Set(s.ctx, key, value, expiration).Err()
}

// LPush 向列表左侧推入元素
func (s *RedisDao) LPush(key string, values interface{}) error {
	return s.client.LPush(s.ctx, key, values).Err()
}

// LLen 获取列表的长度
func (s *RedisDao) LLen(key string) (int64, error) {
	return s.client.LLen(s.ctx, key).Result()
}

// RPop 从列表右侧弹出元素
func (s *RedisDao) RPop(key string) (string, error) {
	return s.client.RPop(s.ctx, key).Result()
}

// SAdd 向集合中添加成员
func (s *RedisDao) SAdd(key string, members ...interface{}) error {
	return s.client.SAdd(s.ctx, key, members...).Err()
}

func (s *RedisDao) SISMMBER(key string, val interface{}, members interface{}) (bool, error) {
	return s.client.SIsMember(s.ctx, key, members).Result()
}

// SMembers 获取集合中的所有成员
func (s *RedisDao) SMembers(key string) ([]string, error) {
	return s.client.SMembers(s.ctx, key).Result()
}

func (s *RedisDao) SREM(key string, members ...interface{}) error {
	return s.client.SRem(s.ctx, key, members...).Err()
}

// HGet 获取哈希表中指定字段的值
func (s *RedisDao) HGet(key, field string) (string, error) {
	return s.client.HGet(s.ctx, key, field).Result()
}

// HGetAll 获取哈希表中的所有字段和值
func (s *RedisDao) HGetAll(key string) (map[string]string, error) {
	return s.client.HGetAll(s.ctx, key).Result()
}

// HSet 设置哈希表中的字段和值
func (s *RedisDao) HSet(key string, values ...interface{}) error {
	return s.client.HSet(s.ctx, key, values...).Err()
}

// HDel 删除哈希表中的字段
func (s *RedisDao) HDel(key string, fields ...string) error {
	return s.client.HDel(s.ctx, key, fields...).Err()
}

// IncrBy 将键的整数值增加指定的增量
func (s *RedisDao) IncrBy(key string, value int64) (int64, error) {
	return s.client.IncrBy(s.ctx, key, value).Result()
}
