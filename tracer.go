package log

import (
	"context"
	"fmt"
	"strings"
	"time"

	opentracing "github.com/opentracing/opentracing-go"
)

type SQLArgs []interface{}

func (args SQLArgs) String() string {
	if len(args) == 0 {
		return "[]"
	}
	var sb strings.Builder
	sb.WriteString("[")
	for i := range args {
		if i > 0 {
			sb.WriteString(",")
		}
		if bs, ok := args[i].([]byte); ok {
			sb.WriteString("`")
			sb.Write(bs)
			sb.WriteString("`")
		} else if s, ok := args[i].(string); ok {
			sb.WriteString("`")
			sb.WriteString(s)
			sb.WriteString("`")
		} else if t, ok := args[i].(time.Time); ok {
			sb.WriteString("`")
			sb.WriteString(t.Format(time.RFC3339Nano))
			sb.WriteString("`")
		} else {
			fmt.Fprint(&sb, args[i])
		}
	}
	sb.WriteString("]")
	return sb.String()
}

func AnyArray(args []interface{}) fmt.Stringer {
	return SQLArgs(args)
}

// SQLTracer 是 github.com/runner-mei/GoBatis 的 Tracer
type SQLTracer struct {
	Logger    Logger
	SpanLevel Level
	LogLevel Level
}

func (w SQLTracer) Write(ctx context.Context, id, sql string, args []interface{}, err error) {
	logger := w.Logger
	if ctx != nil {
		logger = LoggerFromContext(ctx, logger)
		logger = Span(logger, opentracing.SpanFromContext(ctx), w.SpanLevel)
	}
	if w.LogLevel == InfoLevel {
		if err == nil {
			logger.Info(sql, String("id", id), Stringer("args", SQLArgs(args)))
		} else {
			logger.Info(sql, String("id", id), Stringer("args", SQLArgs(args)), Error(err))
		}
		return 
	}
	if err == nil {
		logger.Debug(sql, String("id", id), Stringer("args", SQLArgs(args)))
	} else {
		logger.Debug(sql, String("id", id), Stringer("args", SQLArgs(args)), Error(err))
	}
}

func NewSQLTracer(logger Logger, lvl ...Level) SQLTracer {
	var spanLevel = DefaultSpanLevel
	if len(lvl) > 0 {
		spanLevel = lvl[0]
	}

	return SQLTracer{
		Logger:    logger.AddCallerSkip(4),
		SpanLevel: spanLevel,
		LogLevel:  InfoLevel,
	}
}

func NewDebugSQLTracer(logger Logger, lvl ...Level) SQLTracer {
	var spanLevel = DefaultSpanLevel
	if len(lvl) > 0 {
		spanLevel = lvl[0]
	}

	return SQLTracer{
		Logger:    logger.AddCallerSkip(4),
		SpanLevel: spanLevel,
		LogLevel:  DebugLevel,
	}
}
