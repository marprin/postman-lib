package job

import (
	"context"

	"github.com/opentracing/opentracing-go"
)

//go:generate mockgen -source=./contract.go -destination=./mock/contract.go -package=mock

type (
	JobContract interface {
		Delay(ctx context.Context, jobName string, params JobParam, tracer opentracing.Tracer) (*DelayJobResponse, error)
		DelayIn(ctx context.Context, delayInSec int64, jobName string, params JobParam, tracer opentracing.Tracer) (*DelayJobResponse, error)
	}
)
