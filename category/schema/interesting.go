package schema

import (
	"encoding/json"
	"errors"
	"net"
	"time"
	LOG "watcher/common/log"
)

type whoInteresting struct {
	op         bool
	who        *net.Conn
	accessTime int64
}

func NewWhoInteresting(conn *net.Conn) *whoInteresting{
	return &whoInteresting{op:true, who:conn, accessTime: time.Now().Unix()}
}

func (w *whoInteresting) Update() {
	w.accessTime = time.Now().Unix()
}

func (w *whoInteresting) Close() {
	w.op = false
	(*w.who).Close()
	w.accessTime = 0
}

// true 리턴하면 주어진 seconds을 이미 지났음
func (w *whoInteresting) IsTimeLimit(seconds int64) bool {
	diff := time.Now().Unix() - w.accessTime
	if diff < 0 {
		LOG.Fatal("Diff value is negative")
		return false
	}
	if diff > seconds {
		return true
	}
	return false
}

func (w *whoInteresting) Send(jsonBytes json.RawMessage) error {
	if !w.op {
		return errors.New("empty")
	}
	encoder := json.NewEncoder(*w.who)
	err := encoder.Encode(jsonBytes)
	if err != nil {
		LOG.Error("Encode fail.:", err.Error())
		return err
	}
	return nil
}
