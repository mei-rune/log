package log

import (
	"io"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New(out io.Writer) Logger {
	logger, _ := zap.NewProductionConfig().
		Build(zap.ErrorOutput(zapcore.AddSync(out)))
	return NewLogger(logger)
}
