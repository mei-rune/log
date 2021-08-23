package log

import (
	"strconv"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Field = zapcore.Field
type Level = zapcore.Level
type ArrayMarshaler = zapcore.ArrayMarshaler
type ObjectMarshaler = zapcore.ObjectMarshaler

const (
	// DebugLevel logs are typically voluminous, and are usually disabled in
	// production.
	DebugLevel = zapcore.DebugLevel
	// InfoLevel is the default logging priority.
	InfoLevel = zapcore.InfoLevel
	// WarnLevel logs are more important than Info, but don't need individual
	// human review.
	WarnLevel = zapcore.WarnLevel
	// ErrorLevel logs are high-priority. If an application is running smoothly,
	// it shouldn't generate any error-level logs.
	ErrorLevel = zapcore.ErrorLevel
	// DPanicLevel logs are particularly important errors. In development the
	// logger panics after writing the message.
	DPanicLevel = zapcore.DPanicLevel
	// PanicLevel logs a message, then panics.
	PanicLevel = zapcore.PanicLevel
	// FatalLevel logs a message, then calls os.Exit(1).
	FatalLevel = zapcore.FatalLevel
)

var (

	// Binary constructs a field that carries an opaque binary blob.
	//
	// Binary data is serialized in an encoding-appropriate format. For example,
	// zap's JSON encoder base64-encodes binary blobs. To log UTF-8 encoded text,
	// use ByteString.
	Binary = zap.Binary

	// ByteString constructs a field that carries UTF-8 encoded text as a []byte.
	// To log opaque binary blobs (which aren't necessarily valid UTF-8), use
	// Binary.
	ByteString = zap.ByteString

	// String adds a string-valued key:value pair to a Span.LogFields() record
	String = zap.String

	// Bool adds a bool-valued key:value pair to a Span.LogFields() record
	Bool = zap.Bool

	// Int adds an int-valued key:value pair to a Span.LogFields() record
	Int = zap.Int

	// Int32 adds an int32-valued key:value pair to a Span.LogFields() record
	Int32 = zap.Int32

	// Int64 adds an int64-valued key:value pair to a Span.LogFields() record
	Int64 = zap.Int64

	// Uint32 adds a uint32-valued key:value pair to a Span.LogFields() record
	Uint32 = zap.Uint32

	// Uint64 adds a uint64-valued key:value pair to a Span.LogFields() record
	Uint64 = zap.Uint64

	// Float32 adds a float32-valued key:value pair to a Span.LogFields() record
	Float32 = zap.Float32

	// Float64 adds a float64-valued key:value pair to a Span.LogFields() record
	Float64 = zap.Float64

	// Error adds an error with the key "error" to a Span.LogFields() record
	Error = zap.Error

	// NamedError constructs a field that lazily stores err.Error() under the
	// provided key. Errors which also implement fmt.Formatter (like those produced
	// by github.com/pkg/errors) will also have their verbose representation stored
	// under key+"Verbose". If passed a nil error, the field is a no-op.
	//
	// For the common case in which the key is simply "error", the Error function
	// is shorter and less repetitive.
	NamedError = zap.NamedError

	// Object adds an object-valued key:value pair to a Span.LogFields() record
	Object = zap.Object

	// Noop creates a no-op log field that should be ignored by the tracer.
	Noop = zap.Skip

	// Any adds an any-valued key:value pair to a Span.LogFields() record
	Any = zap.Any

	// Stringer constructs a field with the given key and the output of the value's
	// String method. The Stringer's String method is called lazily.
	Stringer = zap.Stringer

	// Time constructs a Field with the given key and value. The encoder
	// controls how the time is serialized.
	Time = zap.Time

	// Stack constructs a field that stores a stacktrace of the current goroutine
	// under provided key. Keep in mind that taking a stacktrace is eager and
	// expensive (relatively speaking); this function both makes an allocation and
	// takes about two microseconds.
	Stack = zap.Stack

	// Duration constructs a field with the given key and value. The encoder
	// controls how the duration is serialized.
	Duration = zap.Duration
)

type int64Array []int64

func (a int64Array) String() string {
	var sb strings.Builder
	for idx, v := range a {
		if idx != 0 {
			sb.WriteString(",")
		}
		sb.WriteString(strconv.FormatInt(v, 10))
	}
	return sb.String()
}

func Int64Array(name string, a []int64) Field {
	return Stringer(name, int64Array(a))
}

type intArray []int

func (a intArray) String() string {
	var sb strings.Builder
	for idx, v := range a {
		if idx != 0 {
			sb.WriteString(",")
		}
		sb.WriteString(strconv.FormatInt(int64(v), 10))
	}
	return sb.String()
}

func IntArray(name string, a []int) Field {
	return Stringer(name, intArray(a))
}

type uint64Array []uint64

func (a uint64Array) String() string {
	var sb strings.Builder
	for idx, v := range a {
		if idx != 0 {
			sb.WriteString(",")
		}
		sb.WriteString(strconv.FormatUint(v, 10))
	}
	return sb.String()
}

func Uint64Array(name string, a []uint64) Field {
	return Stringer(name, uint64Array(a))
}

type uintArray []uint

func (a uintArray) String() string {
	var sb strings.Builder
	for idx, v := range a {
		if idx != 0 {
			sb.WriteString(",")
		}
		sb.WriteString(strconv.FormatUint(uint64(v), 10))
	}
	return sb.String()
}

func UintArray(name string, a []uint) Field {
	return Stringer(name, uintArray(a))
}

type stringArray []string

func (a stringArray) String() string {
	var sb strings.Builder
	for idx, s := range a {
		if idx != 0 {
			sb.WriteString(",")
		}
		sb.WriteString(s)
	}
	return sb.String()
}

func StringArray(name string, a []string) Field {
	return Stringer(name, stringArray(a))
}

type StringerFunc func() string

func (s StringerFunc) String() string {
	return s()
}

func StringFunc(name string, cb func() string) Field {
	return Stringer(name, StringerFunc(cb))
}