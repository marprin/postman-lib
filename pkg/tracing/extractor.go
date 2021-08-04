package tracing

import (
	"context"

	"github.com/opentracing/opentracing-go"
)

func ExtractTraceID(ctx context.Context, tracer opentracing.Tracer) string {
	span, ctx := opentracing.StartSpanFromContext(ctx, "[tracing][ExtractTraceID]")
	defer span.Finish()

	carrier := map[string]string{}
	_ = tracer.Inject(span.Context(), opentracing.TextMap, opentracing.TextMapCarrier(carrier))
	return carrier["uber-trace-id"]
}

func ExtractHTTP(span opentracing.Span) map[string]string {
	carrier := opentracing.HTTPHeadersCarrier{}
	opentracing.GlobalTracer().Inject(span.Context(), opentracing.HTTPHeaders, carrier)

	result := make(map[string]string)
	for k, v := range carrier {
		result[k] = v[0]
	}

	return result
}
