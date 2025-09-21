package logger

import (
	"loggingLibGo/message"
	"loggingLibGo/sinks"
)

type Config struct {
	MinLevel message.LogLevel
	Sink     sinks.Sink
}
