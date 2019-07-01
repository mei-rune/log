package log

import (
	"context"
	"fmt"

	opentracing "github.com/opentracing/opentracing-go"
)

type sqlArgs []interface{}

func (args sqlArgs) String() string {
	return fmt.Sprintf("%#v", args)
}

// SQLTracer 是 github.com/runner-mei/GoBatis 的 Tracer
type SQLTracer struct {
	Logger    Logger
	SpanLevel Level
}

func (w SQLTracer) Write(ctx context.Context, id, sql string, args []interface{}, err error) {
	logger := w.Logger
	if ctx != nil {
		logger = Span(logger, opentracing.SpanFromContext(ctx), w.SpanLevel)
	}

	if err == nil {
		logger.Info(sql, String("id", id), Stringer("args", sqlArgs(args)))
	} else {
		logger.Info(sql, String("id", id), Stringer("args", sqlArgs(args)), Error(err))
	}
}

func NewSQLTracer(logger Logger, lvl ...Level) SQLTracer {
	var spanLevel = DefaultSpanLevel
	if len(lvl) > 0 {
		spanLevel = lvl[0]
	}

	return SQLTracer{
		Logger:    logger.AddCallerSkip(1),
		SpanLevel: spanLevel,
	}
}
