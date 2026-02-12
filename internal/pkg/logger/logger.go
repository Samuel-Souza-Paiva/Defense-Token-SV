package logger

import (
	"log"
	"os"
)

type Logger struct {
	Prefix string
	out    *log.Logger
}

func (l *Logger) PrintLog(msg string) {
	l.print("LOG", msg, nil)
}

func (l *Logger) PrintSuccess(msg string) {
	l.print("SUCCESS", msg, nil)
}

func (l *Logger) PrintWarning(msg string) {
	l.print("WARN", msg, nil)
}

func (l *Logger) PrintError(msg string, err error) {
	l.print("ERROR", msg, err)
}

func (l *Logger) print(level string, msg string, err error) {
	logger := l.ensure()
	prefix := l.Prefix
	if prefix == "" {
		prefix = "[TS]"
	}
	if err != nil {
		logger.Printf("%s %s %s: %s", prefix, level, msg, err.Error())
		return
	}
	logger.Printf("%s %s %s", prefix, level, msg)
}

func (l *Logger) ensure() *log.Logger {
	if l == nil {
		return log.New(os.Stdout, "", log.LstdFlags)
	}
	if l.out == nil {
		l.out = log.New(os.Stdout, "", log.LstdFlags)
	}
	return l.out
}
