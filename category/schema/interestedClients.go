// 해당 스키마의 데이터에 관심이 있다면, 등록한다.
// 해당 스키마의 데이터가 변경되면 등록된 객체(대상)한테
// 데이터를 송신한다.

package schema

import (
	"encoding/json"
	"errors"
	"net"
	"time"
	LOG "watcher/common/log"
)

type interestedClient struct {
	op         bool
	who        *net.Conn
	accessTime int64
}

func NewInterestedPerson(conn *net.Conn) *interestedClient {
	return &interestedClient{op: true, who:conn, accessTime: time.Now().Unix()}
}

// 시간을 업데이트 한다.
func (w *interestedClient) Update() {
	w.accessTime = time.Now().Unix()
}

// 등록했었던 객체(대상)을 정리한다.
func (w *interestedClient) Release() {
	w.op = false
	(*w.who).Close()
	w.accessTime = 0
}

// true 리턴하면 주어진 seconds을 이미 지났음
func (w *interestedClient) IsTimeLimit(seconds int64) bool {
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

// 데이터를 송신하다. 정해지지 않은 json.RawMessage를 보낸다.
func (w *interestedClient) Send(jsonBytes json.RawMessage) error {
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
