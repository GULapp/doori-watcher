package log

// 만약, 패키지에서 로그함수를 호출하여, 로그를 기록할 경우,
// main패키지가 나중에 호출되므로, 로그파일명을 지정하여, 초기화를 할수가 없다.
// 이를 극복하기 위해, 로그파일명을 초기화 하기전에, 로그함수를 호출한다면, 임의 로그파일에 기록하도록한다.
// 이후 Init 함수가 호출되면, 임의 로그파일에 기록된 내용을 Init함수 호출시, 지정된 로그파일명에,
// 다시 옮기는 작업을 한다.
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
	kFatal   = "FATAL   : "
)

type LEVEL int

const (
	TRACE LEVEL = 0 + iota //0
	DEBUG                  //1
	INFO                   //2
	ERROR                  //3
	FATAL                  //4
)

var (
	instance        *Log
	initCallOnce    bool
	logFullPathFile string
	logfilePerm     os.FileMode
	level           LEVEL
)

type Log struct {
	customLogger *log.Logger
}

func init() {
	logDir := os.Getenv("GUL_LOGPATH")
	if len(logDir) == 0 {
		logDir = "/tmp/"
	}
	procName := filepath.Base(os.Args[0])
	logFullPathFile = logDir + procName + ".log"
	logfilePerm = 0744
	instance = nil
	initCallOnce = false
	level = INFO
	createLogfile(logFullPathFile, logfilePerm)
}

func getInstance() *log.Logger {
	return instance.customLogger
}

// 로그파일을 이름과, 로그파일 생성시, 권한을 셋팅합니다.
// 로그파일명은 절대경로 형식으로 입력 받습니다.
func Init(fileFullPath string, loglevel LEVEL, perm os.FileMode) {
	indicatedLogFile := moveIntoFile(logFullPathFile, fileFullPath, perm)

	logFullPathFile = fileFullPath
	logfilePerm = perm
	level = loglevel
	logger := log.New(indicatedLogFile, "INFO", kLogflag)
	instance = &Log{logger}
}

// previousFile, afterFile 는 FullPath로 입력되어져야 한다.
// previousFile에 있는 contents를 afterFile로 이동한다.
// 이동이 된 후 previousFile 삭제처리된다
// error 발생시, 프로세스 종료
func moveIntoFile(previousFile string, afterFile string, mode os.FileMode) *os.File {
	fileExists(previousFile)

	fileMemory, err := ioutil.ReadFile(previousFile)
	if err != nil {
		panic(err)
		os.Exit(-1)
	} else {
		if err := os.Remove(previousFile); err != nil {
			panic(err)
			os.Exit(-1)
		}
	}

	file, err := os.OpenFile(afterFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, mode)
	if err != nil {
		panic(err)
		os.Exit(-1)
	}

	if _, err = file.Write(fileMemory); err != nil {
		panic(err)
		os.Exit(-1)
	}
	return file
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		panic(err)
		os.Exit(-1)
	}
	//return true
	return !info.IsDir()
}

func createLogfile(filepath string, perm os.FileMode) {
	file, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, perm)
	if err != nil {
		panic(err)
		os.Exit(-1)
	}
	logger := log.New(file, "INFO", kLogflag)
	instance = &Log{customLogger: logger}
}

func Trace(tempPrefix string, format string, v ...interface{}) {
	if level > TRACE {
		return
	}
	getInstance().SetPrefix(fmt.Sprintf("%-8s: ", tempPrefix))
	getInstance().Output(2, fmt.Sprintf(format, v...))
}

func Debug(format string, v ...interface{}) {
	if level > DEBUG {
		return
	}
	getInstance().SetPrefix(kDebug)
	getInstance().Output(2, fmt.Sprintf(format, v...))
}

func Info(format string, v ...interface{}) {
	if level > INFO {
		return
	}
	getInstance().SetPrefix(kInfo)
	getInstance().Output(2, fmt.Sprintf(format, v...))
}

func Error(format string, v ...interface{}) {
	if level > ERROR {
		return
	}
	getInstance().SetPrefix(kError)
	getInstance().Output(2, fmt.Sprintf(format, v...))
}

func Fatal(format string, v ...interface{}) {
	if level > FATAL {
		return
	}
	getInstance().SetPrefix(kFatal)
	getInstance().Output(2, fmt.Sprintf(format, v...))
}
