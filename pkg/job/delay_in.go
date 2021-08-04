package job

import (
	"context"

	"github.com/marprin/postman-lib/pkg/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

func (j *job) DelayIn(ctx context.Context, delayInSec int64, jobName string, params JobParam, tracer opentracing.Tracer) (*DelayJobResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "[Job][DelayIn]")
	defer span.Finish()

	// Inject Trace
	params["uber_trace_id"] = tracing.ExtractTraceID(ctx, tracer)
	span.SetTag("uber_trace_id", params["uber_trace_id"])
	span.SetTag("job.name", jobName)
	span.SetTag("worker.namespace", j.workerNamespace)

	resp, err := j.libWorker.EnqueueIn(jobName, delayInSec, params)
	if err != nil {
		span.SetTag("error", true).LogFields(
			log.String("error delay job", err.Error()),
		)
		return nil, err
	}

	return &DelayJobResponse{
		ID:        resp.ID,
		Name:      resp.Name,
		EnqueueAt: resp.EnqueuedAt,
	}, nil
}
