package logger

import (
	"fmt"
	"log"
	"os"
	"runtime"
)

// Logger is the server logger
type Logger struct {
	logger     *log.Logger
	debug      bool
	trace      bool
	file       bool
	infoLabel  string
	errorLabel string
	fatalLabel string
	debugLabel string
	traceLabel string
}

// NewStdLogger creates a logger with output directed to Stderr
func NewStdLogger(time, file, debug, trace, colors, pid bool) *Logger {
	flags := 0
	if time {
		flags = log.LstdFlags | log.Lmicroseconds
	}

	pre := ""
	if pid {
		pre = pidPrefix()
	}

	l := &Logger{
		logger: log.New(os.Stdout, pre, flags),
		debug:  debug,
		trace:  trace,
		file:   file,
	}

	if colors {
		setColoredLabelFormats(l)
	} else {
		setPlainLabelFormats(l)
	}

	return l
}

// NewFileLogger creates a logger with output directed to a file
func NewFileLogger(filename string, time, file, debug, trace, pid bool) *Logger {
	fileflags := os.O_WRONLY | os.O_APPEND | os.O_CREATE
	f, err := os.OpenFile(filename, fileflags, 0660)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	flags := 0
	if time {
		flags = log.LstdFlags | log.Lmicroseconds
	}

	pre := ""
	if pid {
		pre = pidPrefix()
	}

	l := &Logger{
		logger: log.New(f, pre, flags),
		debug:  debug,
		trace:  trace,
		file:   file,
	}

	setPlainLabelFormats(l)
	return l
}

// Generate the pid prefix string
func pidPrefix() string {
	return fmt.Sprintf("[%d] ", os.Getpid())
}

func setPlainLabelFormats(l *Logger) {
	l.infoLabel = "[INF]"
	l.debugLabel = "[DBG]"
	l.errorLabel = "[ERR]"
	l.fatalLabel = "[FTL]"
	l.traceLabel = "[TRC]"
}

func setColoredLabelFormats(l *Logger) {
	colorFormat := "[\x1b[%dm%s\x1b[0m]"
	l.infoLabel = fmt.Sprintf(colorFormat, 32, "INF")
	l.debugLabel = fmt.Sprintf(colorFormat, 36, "DBG")
	l.errorLabel = fmt.Sprintf(colorFormat, 31, "ERR")
	l.fatalLabel = fmt.Sprintf(colorFormat, 31, "FTL")
	l.traceLabel = fmt.Sprintf(colorFormat, 33, "TRC")
}

// copy from go/src/log/log.go and modify by guojuntao, 2017/09/21
// Cheap integer to fixed-width decimal ASCII.  Give a negative width to avoid zero-padding.
func itoa(i int, wid int) string {
	// Assemble decimal in reverse order.
	var b [20]byte
	bp := len(b) - 1
	for i >= 10 || wid > 1 {
		wid--
		q := i / 10
		b[bp] = byte('0' + i - q*10)
		bp--
		i = q
	}
	// i < 10
	b[bp] = byte('0' + i)
	return string(b[bp:])
}

func (l *Logger) getFileLineFormat(calldepth int) string {
	if l.file {
		_, file, line, ok := runtime.Caller(calldepth)
		if !ok {
			file = "???"
			line = 0
		}
		short := file
		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				short = file[i+1:]
				break
			}
		}
		file = short
		return " " + file + ":" + itoa(line, -1) + ":"
	}
	return ""
}

// Noticef logs a notice statement
func (l *Logger) Noticef(format string, v ...interface{}) {
	l.logger.Printf(l.infoLabel+l.getFileLineFormat(2)+" "+format, v...)
}

func (l *Logger) Noticeln(v ...interface{}) {
	l.logger.Println(append([]interface{}{l.infoLabel + l.getFileLineFormat(2)}, v...)...)
}

// Errorf logs an error statement
func (l *Logger) Errorf(format string, v ...interface{}) {
	l.logger.Printf(l.errorLabel+l.getFileLineFormat(2)+" "+format, v...)
}

func (l *Logger) Errorln(v ...interface{}) {
	l.logger.Println(append([]interface{}{l.errorLabel + l.getFileLineFormat(2)}, v...)...)
}

// Fatalf logs a fatal error
func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.logger.Fatalf(l.fatalLabel+l.getFileLineFormat(2)+" "+format, v...)
}

func (l *Logger) Fatalln(v ...interface{}) {
	l.logger.Fatalln(append([]interface{}{l.fatalLabel + l.getFileLineFormat(2)}, v...)...)
}

// Debugf logs a debug statement
func (l *Logger) Debugf(format string, v ...interface{}) {
	if l.debug {
		l.logger.Printf(l.debugLabel+l.getFileLineFormat(2)+" "+format, v...)
	}
}

func (l *Logger) Debugln(v ...interface{}) {
	if l.debug {
		l.logger.Println(append([]interface{}{l.debugLabel + l.getFileLineFormat(2)}, v...)...)
	}
}

// Tracef logs a trace statement
func (l *Logger) Tracef(format string, v ...interface{}) {
	if l.trace {
		l.logger.Printf(l.traceLabel+l.getFileLineFormat(2)+" "+format, v...)
	}
}

func (l *Logger) Traceln(v ...interface{}) {
	if l.trace {
		l.logger.Println(append([]interface{}{l.traceLabel + l.getFileLineFormat(2)}, v...)...)
	}
}

func (l *Logger) IsDebugEnabled() bool {
	return l.debug
}

func (l *Logger) IsTraceEnabled() bool {
	return l.trace
}
