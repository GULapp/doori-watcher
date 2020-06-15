package common

import (
	"os"
	"testing"
)

var (
	testLog *DrLog
)

func TestNewDrLog(t *testing.T) {
	testLog = NewDrLog("./log", 0744)
}

func TestDrLog_Debug(t *testing.T) {
	testLog.Debug("debug %d\n", 100)
	testLog.Info("info %f\n", 0.0003)
	testLog.FATAL("fatal %s\n", "leejaeseong")
	testLog.Trace("bug1", "check value [%d]\n", 2000)
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
