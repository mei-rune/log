package log

import (
	"io"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New(out io.Writer) Logger {
	outSink := zapcore.Lock(zapcore.AddSync(out))

	cfg := zap.NewProductionConfig()
	logger := zap.New(
		zapcore.NewCore(
			zapcore.NewJSONEncoder(cfg.EncoderConfig),
			outSink,
			cfg.Level,
		),
		zap.ErrorOutput(outSink),
	)
	return NewLogger(logger)
}