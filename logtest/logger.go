package logtest

import (
	"fmt"
	"testing"

	"github.com/runner-mei/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest"
	"go.uber.org/zap/zaptest/observer"
)

type ObservedLogs = observer.ObservedLogs
type LoggedEntry = observer.LoggedEntry

func NewLogger(t testing.TB) log.Logger {
	logger := zaptest.NewLogger(t)
	return log.NewLogger(logger)
}

func New() log.Logger {
	return log.NewStdDefaultLogger()
}

func NewBufferLogger() (log.Logger, *zaptest.Buffer) {
	errSink := &zaptest.Buffer{}
	logger := zap.New(
		zapcore.NewCore(
			zapcore.NewJSONEncoder(zap.NewProductionConfig().EncoderConfig),
			zapcore.Lock(errSink),
			zap.DebugLevel,
		),
		zap.ErrorOutput(errSink),
	)
	return log.NewLogger(logger), errSink
}

func NewObservedLogger() (log.Logger, *ObservedLogs) {
	sf, logs := observer.New(zap.DebugLevel)
	return log.NewLogger(zap.New(sf)), logs
}

func LastEntry(logEntries *ObservedLogs) (bool, LoggedEntry) {
	entries := logEntries.All()
	if len(entries) == 0 {
		return false, LoggedEntry{}
	}
	return true, entries[len(entries)-1]
}

func LastLine(logEntries *zaptest.Buffer) (bool, string) {
	entries := logEntries.Lines()
	if len(entries) == 0 {
		return false, ""
	}
	return true, entries[len(entries)-1]
}

// FieldsMap returns a map for all fields in Context.
func FieldsMap(entries []log.Field) map[string]string {
	encoder := zapcore.NewMapObjectEncoder()
	for _, f := range entries {
		f.AddTo(encoder)
	}

	results := map[string]string{}
	for k, v := range encoder.Fields {
		results[k] = fmt.Sprint(v)
	}
	return results
}
