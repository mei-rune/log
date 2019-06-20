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

	opentracing "github.com/opentracing/opentracing-go"
)

func Span(logger Logger, span opentracing.Span, enabledLevel ...Level) Logger {
	if span == nil {
		return logger
	}

	if len(enabledLevel) > 0 {
		return logger.WithTargets(OutputToTracer(enabledLevel[0], span))
	}

	return logger.WithTargets(OutputToTracer(DefaultSpanLevel, span))
}

func SpanContext(logger Logger, spanContext opentracing.SpanContext, method string, enabledLevel ...Level) (Logger, func()) {
	if spanContext == nil {
		return logger, func() {}
	}

	span := opentracing.StartSpan(method, opentracing.ChildOf(spanContext))
	finish := func() {
		span.Finish()
	}

	if len(enabledLevel) > 0 {
		return logger.WithTargets(OutputToTracer(enabledLevel[0], span)), finish
	}
	return logger.WithTargets(OutputToTracer(DefaultSpanLevel, span)), finish
}

// For returns a context-aware Logger. If the context
// contains an OpenTracing span, all logging calls are also
// echo-ed into the span.
func For(ctx context.Context, args ...interface{}) (Logger, func()) {
	var logger Logger
	var span opentracing.Span
	var spanContext opentracing.SpanContext
	var method string
	var level = DefaultSpanLevel
	var fields []Field

	for _, arg := range args {
		switch value := arg.(type) {
		case Logger:
			logger = value
		case opentracing.Span:
			span = value
		case opentracing.SpanContext:
			spanContext = value
		case string:
			method = value
		case Level:
			level = value
		case Field:
			fields = append(fields, value)
		}
	}

	if logger == nil {
		logger = LoggerOrEmptyFromContext(ctx)
	}
	if len(fields) > 0 {
		logger = logger.With(fields...)
	}

	if span != nil {
		return Span(logger, span, level), noop
	}

	if spanContext != nil {
		return SpanContext(logger, spanContext, method, level)
	}

	if span := opentracing.SpanFromContext(ctx); span != nil {
		return Span(logger, span, level), noop
	}
	return logger, noop
}

var noop = func() {}

func IsEmpty(logger Logger) bool {
	return logger == empty
}
