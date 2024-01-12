package logger

import "context"

type ctxRequestId struct{} // reqId in context

// WithRequestId returns context with requestId value
func WithRequestId(ctx context.Context, requestId string) context.Context {
	return context.WithValue(ctx, ctxRequestId{}, requestId)
}

// GetRequestId returns requestId from context
func GetRequestId(ctx context.Context) (requestId string, ok bool) {
	requestId, ok = ctx.Value(ctxRequestId{}).(string)
	return
}

// WithCtxRequestId returns logger with reqId field from logger
func (sl SLogger) WithCtxRequestId(ctx context.Context) SLogger {
	requestId, ok := GetRequestId(ctx)
	if ok {
		return SLogger{sl.With("requestId", requestId)}
	}
	return sl
}

type ctxUserId struct{} // reqId in context

// WithUserId returns context with userId value
func WithUserId(ctx context.Context, userId string) context.Context {
	return context.WithValue(ctx, ctxUserId{}, userId)
}

// GetUserId returns userId from context
func GetUserId(ctx context.Context) (userId string, ok bool) {
	userId, ok = ctx.Value(ctxUserId{}).(string)
	return
}

// WithCtxUserId returns logger with userId field from logger
func (sl SLogger) WithCtxUserId(ctx context.Context) SLogger {
	userId, ok := GetUserId(ctx)
	if ok {
		return SLogger{sl.With("userId", userId)}
	}
	return sl
}
