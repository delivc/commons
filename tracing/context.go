package tracing

import (
	"context"
	"net/http"
)

type contextKey string

const tracerKey = contextKey("nf-tracer-key")

// WrapWithTracer wraps a context with TracerContext
func WrapWithTracer(r *http.Request, rt *RequestTracer) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), tracerKey, rt))
}

// GetFromContext returns the tracing context from a context.Context
func GetFromContext(ctx context.Context) *RequestTracer {
	val := ctx.Value(tracerKey)
	if val == nil {
		return nil
	}
	entry, ok := val.(*RequestTracer)
	if ok {
		return entry
	}
	return nil
}

// GetTracer shorthand for GetFromContext
func GetTracer(r *http.Request) *RequestTracer {
	return GetFromContext(r.Context())
}
