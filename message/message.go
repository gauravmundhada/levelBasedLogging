package message

import "time"

type LogLevel int

const (
	FATAL LogLevel = iota
	ERROR
	WARN
	INFO
	DEBUG
)

func (l LogLevel) String() string {
	return [...]string{"FATAL", "ERROR", "WARN", "INFO", "DEBUG"}[l]
}

type Message struct {
	Content   string
	Level     LogLevel
	Namespace string
	Timestamp time.Time
}
