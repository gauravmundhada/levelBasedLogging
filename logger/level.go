package logger

type LogLevel int

const (
	FATAL LogLevel = iota //0 (iota starts from 0 and increments by 1 for each line of code in the same const block)
	ERROR                 // 1
	WARN                  // 2
	INFO                  // 3
	DEBUG                 // 4
)

func (l LogLevel) String() string {
	return [...]string{"FATAL", "ERROR", "WARN", "INFO", "DEBUG"}[l]
}
