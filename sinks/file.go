package sinks

import (
	"fmt"
	"loggingLibGo/message"
	"os"
	"path/filepath"
	"time"
)

type FileSink struct {
	FilePath string
	file     *os.File
	logCount int
	maxLogs  int
}

func NewFileSink(filePath string) *FileSink {
	return &FileSink{
		FilePath: filePath,
		maxLogs:  99, // rotate after 100 logs (can change limit as needed)
	}
}

func (f *FileSink) Write(msg message.Message) error {
	// open the file if not opened yet
	if f.file == nil {
		file, err := os.OpenFile(f.FilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		f.file = file
	}

	logEntry := fmt.Sprintf("[%s] [%s] [%s] %s\n",
		msg.Timestamp.Format("2006-01-02 15:04:05"),
		msg.Level,
		msg.Namespace,
		msg.Content,
	)

	// Write to file
	_, err := f.file.WriteString(logEntry)
	if err != nil {
		return err
	}

	f.logCount++

	// check if rotation needed
	if f.logCount > f.maxLogs {
		f.rotateFile()
	}
	return nil
}

// can also implement based on file size or time interval
// can also implement compression for old files
func (f *FileSink) rotateFile() {
	// close the current file
	if f.file != nil {
		f.file.Close()
	}

	// rename the current file
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	ext := filepath.Ext(f.FilePath)
	name := f.FilePath[:len(f.FilePath)-len(ext)]
	newName := fmt.Sprintf("%s_%s%s", timestamp, name, ext) // new name format := timestamp_name.ext

	err := os.Rename(f.FilePath, newName)
	if err != nil {
		fmt.Println("Error rotating file:", err)
		return
	}

	newFile, err := os.OpenFile(f.FilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error creating new log file:", err)
		return
	}

	f.file = newFile
	f.logCount = 0 // reset log count
}
