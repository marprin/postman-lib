package job

import (
	"github.com/gocraft/work"
	"github.com/marprin/postman-lib/pkg/redis"
)

func NewJob(o *Options) JobContract {
	redisPool := redis.NewRedisPool(&redis.Options{
		MaxActive:  o.RedisMaxActive,
		MaxIdle:    o.RedisMaxIdle,
		Timeout:    o.RedisIdleTimeout,
		Connection: o.RedisConnection,
		Password:   o.RedisPassword,
	})

	return &job{
		libWorker:       work.NewEnqueuer(o.WorkerNamespace, redisPool),
		workerNamespace: o.WorkerNamespace,
	}
}
