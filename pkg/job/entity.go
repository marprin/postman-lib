package job

import (
	"time"

	"github.com/gocraft/work"
)

type (
	JobParam work.Q

	DelayJobResponse struct {
		ID        string
		Name      string
		EnqueueAt int64
	}

	Options struct {
		WorkerNamespace  string
		RedisMaxIdle     int
		RedisMaxActive   int
		RedisIdleTimeout time.Duration
		RedisConnection  string
		RedisPassword    string
	}

	job struct {
		libWorker       *work.Enqueuer
		workerNamespace string
	}
)
