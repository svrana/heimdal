package zapped

import (
	"context"

	"go.uber.org/zap"
)

type zappedLogger string

var logger zappedLogger

func NewContext(ctx context.Context, z *zap.SugaredLogger) context.Context {
	return context.WithValue(ctx, logger, z)
}

func FromContext(ctx context.Context) *zap.SugaredLogger {
	if l, ok := ctx.Value(logger).(*zap.SugaredLogger); ok {
		return l
	}
	return zap.NewNop().Sugar()
}
