package mq

func NatsInit(url string) (*NatsPool, error) {
	//  nats conn 初始化连接池
	pool, err := NewDefaultPool(url)
	if err != nil {
		return nil, err
	}

	return pool, nil
}
