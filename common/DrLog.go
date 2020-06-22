package common

import (
	"fmt"
	"log"
	"os"
)

const (
	kLogflag = log.Lshortfile | log.Lmicroseconds
	kDebug   = "DEBUG   : "
	kInfo    = "INFO    : "
	kFatal   = "FATAL   : "
)

var (
	defaultPrefix string
)

type DrLog struct {
	customLogger *log.Logger
}

func NewDrLog(filepath string, perm os.FileMode) *DrLog {
	file, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, perm)
	if err != nil {
		fmt.Errorf("failed to open a logfile")
		os.Exit(-1)
	}
	logger := log.New(file, "INFO", kLogflag)
	return &DrLog{logger}
}

func (l *DrLog) Trace(tempPrefix string, format string, v ...interface{}) {
	l.customLogger.SetPrefix(fmt.Sprintf("%-8s: ", tempPrefix))
	l.customLogger.Output(2, fmt.Sprintf(format, v...))
}

func (l *DrLog) Debug(format string, v ...interface{}) {
	l.customLogger.SetPrefix(kDebug)
	l.customLogger.Output(2, fmt.Sprintf(format, v...))
}

func (l *DrLog) Info(format string, v ...interface{}) {
	l.customLogger.SetPrefix(kInfo)
	l.customLogger.Output(2, fmt.Sprintf(format, v...))
}

func (l *DrLog) FATAL(format string, v ...interface{}) {
	l.customLogger.SetPrefix(kFatal)
	l.customLogger.Output(2, fmt.Sprintf(format, v...))
}
