package gulLog

import (
	"testing"
)

var (
	testLog *Log
)

func TestNewDrLog(t *testing.T) {
	Init("/tmp/leejaeseong.log", 0744)
	Debug("Leejaeseong %s", "leejaeseong")
}
