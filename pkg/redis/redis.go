package redis

import (
	"time"

	"github.com/gomodule/redigo/redis"
)

type (
	Options struct {
		MaxActive  int
		MaxIdle    int
		Timeout    time.Duration
		Connection string
		Password   string
	}
)

func NewRedisPool(o *Options) *redis.Pool {
	return &redis.Pool{
		MaxActive:   o.MaxActive,
		MaxIdle:     o.MaxIdle,
		IdleTimeout: o.Timeout * time.Second,
		Wait:        true,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", o.Connection)
			if err != nil {
				return nil, err
			}

			if o.Password != "" {
				_, err = c.Do("AUTH", o.Password)
				if err != nil {
					return nil, err
				}
			}
			return c, nil
		},
	}
}
