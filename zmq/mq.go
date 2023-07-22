package zmq

import (
	"fmt"
	"github.com/streadway/amqp"
)

type Queue struct {
	conn    *Conn
	channel *amqp.Channel
	queue   *amqp.Queue
}

func (q *Queue) Release() error {
	return pool.Release(q.conn)
}

// Publish
// @Description: 普通发送消息【持久】
// @param name
// @param body
// @return error
func Publish(name QueueName, body []byte) error {
	q, err := queue(name)
	if err != nil {
		return err
	}
	defer q.Release()
	return q.channel.Publish(
		"",
		q.queue.Name,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         body,
		})
}

// Subscribe
// @Description: 监听消息，手动ack
// @param name
// @param sname
// @param ch
// @return error
func Subscribe(name QueueName, sname ServeName) (*Queue, <-chan amqp.Delivery, error) {
	if q, err := queue(name); err != nil {
		return nil, nil, err
	} else {
		del, err := q.channel.Consume(
			q.queue.Name,
			sname.String(),
			false,
			false,
			false,
			false,
			nil,
		)
		return q, del, err
	}
}

// queue
// @Description: 获取队列
// @param name
// @return *amqp.Channel
// @return *amqp.Queue
// @return error
func queue(name QueueName) (*Queue, error) {
	cn, err := pool.Acquire()
	if err != nil {
		return nil, fmt.Errorf("消息队列连接失败 %s", err.Error())
	}
	ch, err := cn.Channel()
	if err != nil {
		_ = pool.Release(cn)
		return nil, fmt.Errorf("消息队列建立通道失败 %s", err.Error())
	}
	qu, err := ch.QueueDeclare(string(name), true, false, false, false, nil)
	if err != nil {
		_ = pool.Release(cn)
		return nil, fmt.Errorf("消息队列建立队列失败 %s", err.Error())
	}
	return &Queue{
		conn:    cn,
		channel: ch,
		queue:   &qu,
	}, nil
}
