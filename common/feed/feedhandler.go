package feed

import (
	"bytes"
	"encoding/gob"
	"net"
	"time"
	LOG "watcher/common/log"
)

// 데이터를 주고 받을 때, 데이터의 총 사이즈의 정보가 필요하다.
// FeedHander는 데이터를 수신 시, 데이터의 총 사이즈를 먼저 읽고,
// 다 읽을 때까지 기다려야 한다.(Timeout)
type DataContainer struct {
	streamLen int
	stream    []byte
}

type FeedHandler struct {
}

func NewFeedHandler() *FeedHandler {
	return &FeedHandler{}
}

func (s *FeedHandler) Send(conn net.Conn, streamBuffer []byte) error {
	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)
	enc.Encode(DataContainer{len(streamBuffer), streamBuffer})

	writedLen, err := conn.Write(buffer.Bytes())
	if writedLen < len(streamBuffer) {
		LOG.Fatal("wrong size")
		return err
	}
	return nil
}

func (s *FeedHandler) Connect(address string) (net.Conn, error) {
	conn, err := net.DialTimeout("tcp", address, 1000*time.Millisecond)
	if err != nil {
		LOG.Fatal("failed to call DialIP ")
		return nil, err
	}
	return conn, nil
}
