package logger

import (
	"loggingLibGo/message"
	"sync"
	"time"
)

var (
	configs []Config
	logChan chan message.Message
	once    sync.Once
	wg      sync.WaitGroup
)

func Configure(c []Config) {
	configs = c      // this line will run every time Configure is called and overwrite previous configs (to avoid overwrite use append)
	once.Do(func() { // this block will run only for the first time Configure is called
		logChan = make(chan message.Message, 200) // buffered channel of size 200 -> 200 messages can be queued
		wg.Add(1)                                 // will wait for one goroutine
		go worker()
	})
}

func Log(msg message.Message) {
	log := message.Message{
		Content:   msg.Content,
		Level:     msg.Level,
		Namespace: msg.Namespace,
		Timestamp: time.Now(),
	}
	select {
	case logChan <- log:
		// message successfully sent to channel (can be handled if needed)
	default:
		// channel is full, drop the message (can handle this case if needed)
		// e.x - can implement log rotation
	}
}

func worker() {
	defer wg.Done() // signal that the goroutine is done
	for msg := range logChan {
		for _, cfg := range configs {
			if msg.Level <= cfg.MinLevel {
				cfg.Sink.Write(msg)
			}
		}
	}
}

func Close() {
	close(logChan) // close the channel to signal the worker to stop
	wg.Wait()      // wait for the worker to finish (process all messages)
}
