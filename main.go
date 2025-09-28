package main

import (
	"fmt"
	"loggingLibGo/logger"
	"loggingLibGo/message"
	"loggingLibGo/sinks"
)

func main() {
	elasticSink := sinks.NewElasticSink("http://localhost:9200", "logs")

	logger.Configure([]logger.Config{ // example with multiple sinks(can add different levels and sinks)
		{
			MinLevel: message.DEBUG,
			Sink:     sinks.NewConsoleSink(),
		},
		{
			MinLevel: message.INFO,
			Sink:     sinks.NewFileSink("app.log"),
		},
		{
			MinLevel: message.DEBUG,
			Sink:     elasticSink,
		},
	})

	logger.Log(message.Message{
		Content:   "Application started",
		Level:     message.DEBUG,
		Namespace: "main",
	})

	logger.Log(message.Message{
		Content:   "Test log to Elastic",
		Level:     message.INFO,
		Namespace: "elastic",
	})

	//Simulate multiple log entries for testing rotation
	for range 101 {
		logger.Log(message.Message{
			Content:   "logging in the file",
			Level:     message.INFO,
			Namespace: "main",
		})
	}

	logs, err := elasticSink.Search(sinks.SearchQuery{
		Level:     "INFO",
		Namespace: "elastic",
		StartDate: "",
		EndDate:   "",
	})

	if err != nil {
		fmt.Println("Search error:", err)
		return
	}

	for _, log := range logs {
		fmt.Println("Found log:", log)
	}

	logger.Close()

}
