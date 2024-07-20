package mq

import "github.com/nats-io/go-nats"

/*
基于队列组的主题订阅：
具有相同队列名称的所有订阅都将形成一个队列组。使用队列语义，每个消息将仅传递给每个队列组的一个订阅服务器。
您可以拥有任意数量的队列组。普通订阅服务器将继续按预期工作。
*/
func (self *NatsPool) QueueSubscribe(subject, queue string, handler nats.Handler) error {
	conn, err := self.Get()
	if err != nil {
		return err
	}
	defer self.Put(conn)

	encodedConn, err := nats.NewEncodedConn(conn, nats.DEFAULT_ENCODER)
	if err != nil {
		return err
	}

	_, err = encodedConn.QueueSubscribe(subject, queue, handler)
	return err
}
