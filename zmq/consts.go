package zmq

type QueueName string

type ServeName string

func (n ServeName) String() string {
	return string(n)
}
