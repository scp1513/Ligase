package logger

var logger *Logger

func InitLogger(time, file, debug, trace, colors, pid bool) {
	logger = NewStdLogger(time, file, debug, trace, colors, pid)
}

func GetLogger() *Logger {
	return logger
}
