package common

// 만약, 패키지에서 로그함수를 호출하여, 로그를 기록할 경우,
// 지정된 로그파일

import (
	"fmt"
	"io/ioutil"
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
	logfilePerm os.FileMode
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
	logfilePerm = 0744
	instance = nil
}

func getInstance() *log.Logger{
	if instance == nil {
		Init(logfileName, logfilePerm)
	}
	return instance.customLogger
}

func Init(filepath string, perm os.FileMode) {
	if logfileName != filepath {
		fileMemory, err := ioutil.ReadFile(logfileName)
		if err != nil {
			createLogfile(filepath, perm)
			return
		}

		logfileName = filepath
		logfilePerm = perm
		instance = nil

		createLogfile(filepath, perm)

		// move temporarily log file to appointed log file
		instance.customLogger.Print(string(fileMemory))
	} else{
		createLogfile(filepath, perm)
	}
}

func createLogfile(filepath string, perm os.FileMode){
	file, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, perm)
	if err != nil {
		fmt.Errorf("failed to open a logfile")
		os.Exit(-1)
	}
	logger := log.New(file, "INFO", kLogflag)
	instance = &Log{customLogger:logger}
}

func Trace(tempPrefix string, format string, v ...interface{}) {
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
