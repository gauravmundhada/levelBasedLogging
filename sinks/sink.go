package sinks

import "loggingLibGo/message"

type Sink interface {
	Write(msg message.Message) error
}
