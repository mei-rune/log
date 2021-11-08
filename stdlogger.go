package log

import (
	"fmt"
	stdlog "log"

	"go.uber.org/zap"
)

// stdlogger delegates all calls to the underlying zap.Logger
type stdlogger struct {
	callerSkip int
	name       string
	fields     []Field

	logger *stdlog.Logger
}

func (l stdlogger) Sync() error {
	return nil
}

func (l stdlogger) ToStdLogger() *stdlog.Logger {
	return l.logger
}

func (l stdlogger) log(level Level, msg string, fields []Field) {
	if l.logger == nil {
		stdlog.Output(l.callerSkip, fmt.Sprint(level, " ", msg, " ", append(l.fields, fields...)))
	} else {
		l.logger.Output(l.callerSkip, fmt.Sprint(level, " ", msg, " ", append(l.fields, fields...)))
	}

	if level == FatalLevel {
		panic(msg)
	}
}

func (l stdlogger) logf(level Level, msgfmt string, args []interface{}) {
	msg := fmt.Sprintf(msgfmt, args...)
	if l.logger == nil {
		stdlog.Output(l.callerSkip, level.String()+" "+msg+" "+fmt.Sprint(l.fields))
	} else {
		l.logger.Output(l.callerSkip, level.String()+" "+msg+" "+fmt.Sprint(l.fields))
	}

	if level == FatalLevel {
		panic(msg)
	}
}

// Panic logs an panic msg with fields and panic
func (l stdlogger) Panic(msg string, fields ...Field) {
	l.log(PanicLevel, msg, fields)
}

// Debug logs an debug msg with fields
func (l stdlogger) Debug(msg string, fields ...Field) {
	l.log(DebugLevel, msg, fields)
}

// Info logs an info msg with fields
func (l stdlogger) Info(msg string, fields ...Field) {
	l.log(InfoLevel, msg, fields)
}

// Warn logs an error msg with fields
func (l stdlogger) Warn(msg string, fields ...Field) {
	l.log(WarnLevel, msg, fields)
}

// Error logs an error msg with fields
func (l stdlogger) Error(msg string, fields ...Field) {
	l.log(ErrorLevel, msg, fields)
}

// Fatal logs a fatal error msg with fields
func (l stdlogger) Fatal(msg string, fields ...Field) {
	l.log(FatalLevel, msg, fields)
}

// Debugw logs an debug msg with fields
func (l stdlogger) Debugw(msg string, keyAndValues ...interface{}) {
	fields := SweetenFields(l, keyAndValues)
	l.log(DebugLevel, msg, fields)
}

// Infow logs an info msg with fields
func (l stdlogger) Infow(msg string, keyAndValues ...interface{}) {
	fields := SweetenFields(l, keyAndValues)
	l.log(InfoLevel, msg, fields)
}

// Warnw logs an error msg with fields
func (l stdlogger) Warnw(msg string, keyAndValues ...interface{}) {
	fields := SweetenFields(l, keyAndValues)
	l.log(WarnLevel, msg, fields)
}

// Errorw logs an error msg with fields
func (l stdlogger) Errorw(msg string, keyAndValues ...interface{}) {
	fields := SweetenFields(l, keyAndValues)
	l.log(ErrorLevel, msg, fields)
}

// Fatalw logs a fatal error msg with fields
func (l stdlogger) Fatalw(msg string, keyAndValues ...interface{}) {
	fields := SweetenFields(l, keyAndValues)
	l.log(FatalLevel, msg, fields)
}

// Debugf logs an debug msg with arguments
func (l stdlogger) Debugf(msg string, args ...interface{}) {
	l.logf(FatalLevel, msg, args)
}

// Infow logs an info msg with arguments
func (l stdlogger) Infof(msg string, args ...interface{}) {
	l.logf(InfoLevel, msg, args)
}

// Warnw logs an error msg with arguments
func (l stdlogger) Warnf(msg string, args ...interface{}) {
	l.logf(WarnLevel, msg, args)
}

// Errorw logs an error msg with arguments
func (l stdlogger) Errorf(msg string, args ...interface{}) {
	l.logf(ErrorLevel, msg, args)
}

// Fatalw logs a fatal error msg with arguments
func (l stdlogger) Fatalf(msg string, args ...interface{}) {
	l.logf(FatalLevel, msg, args)
}

// With creates a child logger, and optionally adds some context fields to that logger.
func (l stdlogger) With(fields ...Field) Logger {
	return stdlogger{fields: append(l.fields, fields...), logger: l.logger, callerSkip: l.callerSkip}
}

// With creates a child logger, and optionally adds some context fields to that logger.
func (l stdlogger) WithTargets(targets ...Target) Logger {
	if len(targets) == 0 {
		return l
	}
	return appendLogger{logger: l.AddCallerSkip(1),
		target: Tee(targets)}
}

// Named adds a new path segment to the logger's name. Segments are joined by
// periods. By default, Loggers are unnamed.
func (l stdlogger) Named(name string) Logger {
	newName := l.name
	if newName == "" {
		newName = name
	} else {
		newName = newName + "." + name
	}
	return stdlogger{name: newName, logger: l.logger, callerSkip: l.callerSkip}
}

func (l stdlogger) AddCallerSkip(level int) Logger {
	return stdlogger{name: l.name, fields: l.fields, logger: l.logger, callerSkip: l.callerSkip + level}
}

func (l stdlogger) Unwrap() *zap.Logger {
	return nil
}

func NewStdLogger(logger *stdlog.Logger) Logger {
	return stdlogger{callerSkip: 3, logger: logger}
}

func NewStdDefaultLogger() Logger {
	return stdlogger{callerSkip: 3}
}
