package feed

import (
	LOG "go_monitoring/common/log"
	"net"
	"time"
)

type FeedHandler struct {
}

func NewFeedHandler() *FeedHandler {
	return &FeedHandler{}
}

func (s *FeedHandler) Send(conn net.Conn, streamBuffer []byte) error {
	writedLen,err :=conn.Write(streamBuffer)
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
