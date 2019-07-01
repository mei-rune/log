// Copyright (c) 2017 Uber Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package log

import (
	"context"
	"log"

	"go.uber.org/zap"
)

// Logger is a simplified abstraction of the zap.Logger
type Logger interface {
	Sync() error
	Panic(msg string, fields ...Field)
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Fatal(msg string, fields ...Field)

	Debugw(msg string, fields ...interface{})
	Infow(msg string, fields ...interface{})
	Errorw(msg string, fields ...interface{})
	Warnw(msg string, fields ...interface{})
	Fatalw(msg string, fields ...interface{})

	Debugf(msg string, values ...interface{})
	Infof(msg string, values ...interface{})
	Errorf(msg string, values ...interface{})
	Warnf(msg string, values ...interface{})
	Fatalf(msg string, values ...interface{})

	AddCallerSkip(int) Logger
	With(fields ...Field) Logger
	WithTargets(targets ...Target) Logger
	Named(name string) Logger
	Unwrap() *zap.Logger
	ToStdLogger() *log.Logger
}

// zaplogger delegates all calls to the underlying zap.Logger
type zaplogger struct {
	logger  *zap.Logger
	sugared *zap.SugaredLogger
}

func (l zaplogger) Sync() error {
	return l.logger.Sync()
}

func (l zaplogger) ToStdLogger() *log.Logger {
	return zap.NewStdLog(l.logger)
}

// Panic logs an panic msg with fields and panic
func (l zaplogger) Panic(msg string, fields ...Field) {
	l.logger.Panic(msg, fields...)
}

// Debug logs an debug msg with fields
func (l zaplogger) Debug(msg string, fields ...Field) {
	l.logger.Debug(msg, fields...)
}

// Info logs an info msg with fields
func (l zaplogger) Info(msg string, fields ...Field) {
	l.logger.Info(msg, fields...)
}

// Warn logs an error msg with fields
func (l zaplogger) Warn(msg string, fields ...Field) {
	l.logger.Warn(msg, fields...)
}

// Error logs an error msg with fields
func (l zaplogger) Error(msg string, fields ...Field) {
	l.logger.Error(msg, fields...)
}

// Fatal logs a fatal error msg with fields
func (l zaplogger) Fatal(msg string, fields ...Field) {
	l.logger.Fatal(msg, fields...)
}

// Debugw logs an debug msg with fields
func (l zaplogger) Debugw(msg string, fields ...interface{}) {
	l.sugared.Debugw(msg, fields...)
}

// Infow logs an info msg with fields
func (l zaplogger) Infow(msg string, fields ...interface{}) {
	l.sugared.Infow(msg, fields...)
}

// Warnw logs an error msg with fields
func (l zaplogger) Warnw(msg string, fields ...interface{}) {
	l.sugared.Warnw(msg, fields...)
}

// Errorw logs an error msg with fields
func (l zaplogger) Errorw(msg string, fields ...interface{}) {
	l.sugared.Errorw(msg, fields...)
}

// Fatalw logs a fatal error msg with fields
func (l zaplogger) Fatalw(msg string, fields ...interface{}) {
	l.sugared.Fatalw(msg, fields...)
}

// Debugf logs an debug msg with arguments
func (l zaplogger) Debugf(msg string, args ...interface{}) {
	l.sugared.Infof(msg, args...)
}

// Infow logs an info msg with arguments
func (l zaplogger) Infof(msg string, args ...interface{}) {
	l.sugared.Infof(msg, args...)
}

// Warnw logs an error msg with arguments
func (l zaplogger) Warnf(msg string, args ...interface{}) {
	l.sugared.Warnf(msg, args...)
}

// Errorw logs an error msg with arguments
func (l zaplogger) Errorf(msg string, args ...interface{}) {
	l.sugared.Errorf(msg, args...)
}

// Fatalw logs a fatal error msg with arguments
func (l zaplogger) Fatalf(msg string, args ...interface{}) {
	l.sugared.Fatalf(msg, args...)
}

// With creates a child logger, and optionally adds some context fields to that logger.
func (l zaplogger) With(fields ...Field) Logger {
	newL := l.logger.With(fields...)
	return zaplogger{logger: newL, sugared: newL.Sugar()}
}

// With creates a child logger, and optionally adds some context fields to that logger.
func (l zaplogger) WithTargets(targets ...Target) Logger {
	if len(targets) == 0 {
		return l
	}
	newL := l.logger.WithOptions(zap.AddCallerSkip(1))
	return appendLogger{logger: zaplogger{logger: newL, sugared: newL.Sugar()},
		target: Tee(targets)}
}

// Named adds a new path segment to the logger's name. Segments are joined by
// periods. By default, Loggers are unnamed.
func (l zaplogger) Named(name string) Logger {
	newL := l.logger.Named(name)
	return zaplogger{logger: newL, sugared: newL.Sugar()}
}

func (l zaplogger) AddCallerSkip(level int) Logger {
	logger := l.logger.WithOptions(zap.AddCallerSkip(level))
	return zaplogger{logger: logger, sugared: logger.Sugar()}
}

func (l zaplogger) Unwrap() *zap.Logger {
	return l.logger
}

func NewLogger(logger *zap.Logger) Logger {
	logger = logger.WithOptions(zap.AddCallerSkip(1))
	return zaplogger{logger: logger, sugared: logger.Sugar()}
}

// Logger is a simplified abstraction of the zap.Logger
type emptyLogger struct{}

func (empty emptyLogger) Sync() error { return nil }
func (empty emptyLogger) ToStdLogger() *log.Logger {
	return nil
}
func (empty emptyLogger) Panic(msg string, fields ...Field) {
	panic(msg)
}
func (empty emptyLogger) Debug(msg string, fields ...Field) {}
func (empty emptyLogger) Info(msg string, fields ...Field)  {}
func (empty emptyLogger) Error(msg string, fields ...Field) {}
func (empty emptyLogger) Warn(msg string, fields ...Field)  {}
func (empty emptyLogger) Fatal(msg string, fields ...Field) {}

func (empty emptyLogger) Debugw(msg string, fields ...interface{}) {}
func (empty emptyLogger) Infow(msg string, fields ...interface{})  {}
func (empty emptyLogger) Errorw(msg string, fields ...interface{}) {}
func (empty emptyLogger) Warnw(msg string, fields ...interface{})  {}
func (empty emptyLogger) Fatalw(msg string, fields ...interface{}) {}

func (empty emptyLogger) Debugf(msg string, values ...interface{}) {}
func (empty emptyLogger) Infof(msg string, values ...interface{})  {}
func (empty emptyLogger) Errorf(msg string, values ...interface{}) {}
func (empty emptyLogger) Warnf(msg string, values ...interface{})  {}
func (empty emptyLogger) Fatalf(msg string, values ...interface{}) {}

func (empty emptyLogger) AddCallerSkip(level int) Logger { return empty }
func (empty emptyLogger) With(fields ...Field) Logger    { return empty }
func (empty emptyLogger) Named(name string) Logger       { return empty }
func (empty emptyLogger) WithTargets(targets ...Target) Logger {
	if len(targets) == 0 {
		return empty
	}
	return appendLogger{logger: empty, target: Tee(targets)}
}
func (empty emptyLogger) Unwrap() *zap.Logger { return nil }

// Empty a nil logger
var empty Logger = emptyLogger{}

func Empty() Logger {
	return empty
}

type loggerKey struct{}

var activeLoggerKey = loggerKey{}

// ContextWithLogger returns a new `context.Context` that holds a reference to
// `logger`'s LoggerContext.
func ContextWithLogger(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, activeLoggerKey, logger)
}

// LoggerFromContext returns the `logger` previously associated with `ctx`, or
// `nil` if no such `logger` could be found.
func LoggerFromContext(ctx context.Context) Logger {
	val := ctx.Value(activeLoggerKey)
	if sp, ok := val.(Logger); ok {
		return sp
	}
	return nil
}

// LoggerOrEmptyFromContext returns the `logger` previously associated with `ctx`, or
// `Empty` if no such `logger` could be found.
func LoggerOrEmptyFromContext(ctx context.Context) Logger {
	val := ctx.Value(activeLoggerKey)
	if sp, ok := val.(Logger); ok {
		return sp
	}
	return empty
}
