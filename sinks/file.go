package sinks

import (
	"fmt"
	"loggingLibGo/message"
	"os"
)

type FileSink struct {
	FilePath string
}

func NewFileSink(filePath string) *FileSink {
	return &FileSink{FilePath: filePath}
}

func (f *FileSink) Write(msg message.Message) error {
	file, err := os.OpenFile(f.FilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close() // will execute at last

	line := fmt.Sprintf(
		"[%s] Level - [%s] Namespace - [%s] LogDetaisl - [%s]\n",
		msg.Timestamp.Format("2006-01-02 15:04:05"),
		msg.Level,
		msg.Namespace,
		msg.Content,
	)

	_, err = file.WriteString(line)
	return err
}
