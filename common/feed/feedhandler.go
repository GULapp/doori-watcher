package feed

import (
	"encoding/json"
	"net"
	"time"
	LOG "watcher/common/log"
)

// 데이터를 보낼 때, 데이터의 형식의 이름과, 데이터타입을 그대로 보내기 위한
// 소켓 컨테이너 구조체다. DataObject 는 구조체로, Json tag가 존재해야 한다.
type DataContainer struct {
	ObjectName string          `json:"name"`
	ObjectData json.RawMessage `json:"data"`
}

type FeedHandler struct {
}

func NewFeedHandler() *FeedHandler {
	return &FeedHandler{}
}

// conn객체를 이용해서 streamBuffer에 저장된 바이트 데이터를 보냅니다.
func (s *FeedHandler) Send(conn net.Conn, jsonBytes DataContainer) error {
	encoder := json.NewEncoder(conn)

	err := encoder.Encode(jsonBytes)
	if err != nil {
		LOG.Error("Encode fail.:", err.Error())
		return err
	}
	return nil
}

// address에 tcp 연결합니다.
func (s *FeedHandler) Connect(address string) (net.Conn, error) {
	conn, err := net.DialTimeout("tcp", address, 1000*time.Millisecond)
	if err != nil {
		LOG.Fatal("failed to call DialIP ")
		return nil, err
	}
	return conn, nil
}
