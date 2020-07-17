package common

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

const (
	kLogflag = log.Lshortfile | log.Lmicroseconds
	kDebug   = "DEBUG   : "
	kInfo    = "INFO    : "
	kError   = "ERROR   : "
	kFatal   = "Fatal   : "
)

var (
	instance *Log
	logfileName string
)

type Log struct {
	customLogger *log.Logger
}

func init(){
	logDir := os.Getenv("GUL_LOGPATH")
	if len(logDir) == 0 {
		logDir = "/tmp/"
	}
	procName := filepath.Base(os.Args[0])
	logfileName = logDir+procName+".log"
	instance = Init(logfileName, 0744)
}

func getInstance() *log.Logger{
	if instance == nil {
		instance = Init(logfileName, 0744)
	}
	return instance.customLogger
}

func Init(filepath string, perm os.FileMode) *Log {
	file, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, perm)
	if err != nil {
		fmt.Errorf("failed to open a logfile")
		os.Exit(-1)
	}
	logger := log.New(file, "INFO", kLogflag)
	return &Log{logger}
}

func  Trace(tempPrefix string, format string, v ...interface{}) {
	getInstance().SetPrefix(fmt.Sprintf("%-8s: ", tempPrefix))
	getInstance().Output(2, fmt.Sprintf(format, v...))
}

func Debug(format string, v ...interface{}) {
	getInstance().SetPrefix(kDebug)
	getInstance().Output(2, fmt.Sprintf(format, v...))
}

func Info(format string, v ...interface{}) {
	getInstance().SetPrefix(kInfo)
	getInstance().Output(2, fmt.Sprintf(format, v...))
}

func Error(format string, v ...interface{}) {
	getInstance().SetPrefix(kError)
	getInstance().Output(2, fmt.Sprintf(format, v...))
}

func Fatal(format string, v ...interface{}) {
	getInstance().SetPrefix(kFatal)
	getInstance().Output(2, fmt.Sprintf(format, v...))
}
