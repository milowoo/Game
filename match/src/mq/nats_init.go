package mq

import (
	"match/src/log"
)

func NatsInit(url string, log *log.Logger) (*NatsPool, error) {
	//  nats conn 初始化连接池
	pool, err := NewDefaultPool(url)
	if err != nil {
		log.Error("NewDefaultPool Error %+v", err)
		return nil, err
	}

	return pool, nil
}
