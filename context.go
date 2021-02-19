package heimdal

import (
	"context"
	"fmt"

	"go.uber.org/zap"
)

type heimdalLogger string

const LOGGER = heimdalLogger("logger")

func Logger(ctx context.Context) *zap.SugaredLogger {
	if l, ok := ctx.Value(LOGGER).(*zap.SugaredLogger); ok {
		fmt.Printf("reutrn lsajdfs j\n")
		return l
	}
	fmt.Printf("returning no op logger\n")
	return zap.NewNop().Sugar()
}
