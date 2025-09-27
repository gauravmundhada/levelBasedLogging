package main

import (
	"loggingLibGo/logger"
	"loggingLibGo/message"
	"loggingLibGo/sinks"
)

func main() {
	logger.Configure([]logger.Config{ // example with multiple sinks(can add different levels and sinks)
		{
			MinLevel: message.DEBUG,
			Sink:     sinks.NewConsoleSink(),
		},
		{
			MinLevel: message.INFO,
			Sink:     sinks.NewFileSink("app.log"),
		},
	})

	logger.Log(message.Message{
		Content:   "Application started",
		Level:     message.DEBUG,
		Namespace: "main",
	})

	// Simulate multiple log entries for testing rotation
	for i := 0; i < 101; i++ {
		logger.Log(message.Message{
			Content:   "logging in the file",
			Level:     message.INFO,
			Namespace: "main",
		})
	}

	logger.Close()

}
