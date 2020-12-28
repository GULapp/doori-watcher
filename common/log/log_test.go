package log

import (
	"testing"
)

var (
	testLog *Log
)

func TestNewDrLog(t *testing.T) {
	Init("/tmp/leejaeseong.log", DEBUG, 0744)
	Debug("Leejaeseong %s", "leejaeseong")
}
