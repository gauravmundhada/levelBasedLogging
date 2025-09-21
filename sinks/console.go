package sinks

import (
	"fmt"
	"loggingLibGo/message"
)

type ConsoleSink struct{}

func NewConsoleSink() *ConsoleSink {
	return &ConsoleSink{}
}

func (c *ConsoleSink) Write(msg message.Message) error {
	fmt.Printf("[%s] Level - [%s] Namespace - [%s] LogDetails - [%s]\n",
		msg.Timestamp.Format("2006-01-02 15:04:05"),
		msg.Level,
		msg.Namespace,
		msg.Content,
	)
	return nil
}
