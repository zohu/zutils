package zmq

import (
	"fmt"
	"github.com/streadway/amqp"
)

func Acquire() (*Conn, error) {
	return pool.Acquire()
}

func Release(c *Conn) error {
	return pool.Release(c)
}

// Publish
// @Description: 普通发送消息【持久】
// @param name
// @param body
// @return error
func Publish(name QueueName, body []byte) error {
	if c, q, err := queue(name); err != nil {
		return err
	} else {
		return c.Publish(
			"",
			q.Name,
			false,
			false,
			amqp.Publishing{
				DeliveryMode: amqp.Persistent,
				ContentType:  "application/json",
				Body:         body,
			})
	}
}

// Subscribe
// @Description: 监听消息，手动ack
// @param name
// @param sname
// @param ch
// @return error
func Subscribe(name QueueName, sname ServeName) (<-chan amqp.Delivery, error) {
	if c, q, err := queue(name); err != nil {
		return nil, err
	} else {
		return c.Consume(
			q.Name,
			sname.String(),
			false,
			false,
			false,
			false,
			nil,
		)
	}
}

// queue
// @Description: 获取队列
// @param name
// @return *amqp.Channel
// @return *amqp.Queue
// @return error
func queue(name QueueName) (*amqp.Channel, *amqp.Queue, error) {
	cn, err := Acquire()
	if err != nil {
		return nil, nil, fmt.Errorf("消息队列连接失败 %s", err.Error())
	}
	ch, err := cn.Channel()
	if err != nil {
		return nil, nil, fmt.Errorf("消息队列建立通道失败 %s", err.Error())
	}
	queue, err := ch.QueueDeclare(string(name), true, false, false, false, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("消息队列建立队列失败 %s", err.Error())
	}
	return ch, &queue, nil
}
