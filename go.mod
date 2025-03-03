module github.com/runner-mei/log

go 1.13

require (
	github.com/opentracing/opentracing-go v1.2.0
	github.com/stretchr/testify v1.4.0
	go.uber.org/multierr v1.6.0
	go.uber.org/zap v1.16.0
	golang.org/x/exp v0.0.0-00010101000000-000000000000
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
)

replace golang.org/x/exp => github.com/mei-rune/golang_exp_for_go120 v0.0.0-20250303053821-1e7433e4f2f2
