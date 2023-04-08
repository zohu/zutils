package zmq

import (
	"context"
	"fmt"
	"github.com/streadway/amqp"
	"github.com/zohu/zlog"
	"github.com/zohu/zutils"
	"sync"
	"time"
)

type RabbitmqConfig struct {
	Host     string `yaml:"host" mapstructure:"host"`
	Username string `yaml:"username" mapstructure:"username"`
	Password string `yaml:"password" mapstructure:"password"`
	MinAlive int    `yaml:"min_alive" mapstructure:"min_alive"`
	MaxAlive int    `yaml:"max_alive" mapstructure:"max_alive"`
	MaxIdle  int    `yaml:"max_idle" mapstructure:"max_idle"`
}

type options struct {
	ctx              context.Context
	url              string
	minAlive         int
	maxAlive         int
	maxIdle          int
	idleCloseTimeOut int64
}

type Stats struct {
	total   int
	created int64
}

type Conn struct {
	*amqp.Connection
	id        string
	lastUseAt int64
}

type Pool struct {
	opts    *options
	synchro struct {
		sync.Mutex
		conn  chan *Conn
		stats Stats
	}
	factory func() (*Conn, error)
}

var pool *Pool

func Init(conf *RabbitmqConfig) {
	ops := &options{
		ctx:              context.Background(),
		url:              fmt.Sprintf("amqp://%s:%s@%s/", conf.Username, conf.Password, conf.Host),
		minAlive:         conf.MinAlive,
		maxAlive:         conf.MaxAlive,
		maxIdle:          conf.MaxIdle,
		idleCloseTimeOut: 600,
	}
	pool = &Pool{
		opts: ops,
		factory: func() (*Conn, error) {
			c, err := amqp.Dial(ops.url)
			if err != nil {
				return nil, err
			}
			return &Conn{
				c,
				zutils.NewUuid(),
				time.Now().Unix(),
			}, nil
		},
	}
	pool.synchro.conn = make(chan *Conn, pool.opts.maxAlive)
	if pool.opts.minAlive > 0 {
		pool.initConnection()
	}

	go pool.timerConnectionFactory()
}

func (p *Pool) Stats() *Stats {
	stats := new(Stats)
	*stats = p.synchro.stats
	return stats
}

func (p *Pool) Acquire() (*Conn, error) {
	var c *Conn
	var err error
	for {
		c, err = p.getOrCreate()
		if err != nil {
			return nil, err
		}
		if c.IsClosed() {
			if err = p.Destroy(); err != nil {
				return nil, err
			}
		} else {
			return c, nil
		}
	}
}

func (p *Pool) getOrCreate() (*Conn, error) {
	select {
	case connection := <-p.synchro.conn:
		connection.lastUseAt = time.Now().Unix()
		return connection, nil
	default:
	}
	p.synchro.Lock()
	if p.synchro.stats.total >= p.opts.maxAlive {
		p.synchro.Unlock()
		connection := <-p.synchro.conn
		connection.lastUseAt = time.Now().Unix()
		return connection, nil
	}
	connection, err := p.factory()
	if err != nil {
		p.synchro.Unlock()
		return nil, err
	}
	p.synchro.stats.total++
	p.synchro.stats.created++
	p.synchro.Unlock()
	return connection, nil
}

func (p *Pool) Release(con *Conn) error {
	p.synchro.Lock()
	p.synchro.conn <- con
	p.synchro.Unlock()
	return nil
}

func (p *Pool) Destroy() error {
	p.synchro.Lock()
	p.synchro.stats.total--
	p.synchro.Unlock()
	return nil
}

func (p *Pool) initConnection() {
	p.synchro.Lock()
	defer p.synchro.Unlock()
	for i := 0; i < p.opts.minAlive; i++ {
		con, err := p.factory()
		if err != nil {
			zlog.Errorf("初始化最小连接失败", err.Error())
			continue
		}
		p.synchro.conn <- con
		p.synchro.stats.total++
		p.synchro.stats.created++
	}
}

func (p *Pool) timerConnectionFactory() {
	ticker := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-ticker.C:
			p.synchro.Lock()
			// 释放超时等待连接
			if p.synchro.stats.total > p.opts.maxIdle {
				select {
				case connection := <-p.synchro.conn:
					if time.Now().Unix()-connection.lastUseAt > p.opts.idleCloseTimeOut {
						p.synchro.stats.total--
					} else {
						p.synchro.conn <- connection
					}
				default:

				}
			}
			p.synchro.Unlock()
		case <-p.opts.ctx.Done():
			return
		}
	}
}
