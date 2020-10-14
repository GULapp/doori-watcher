package feed

import (
	"encoding/json"
	"net"
	"os"
	"sync"
	LOG "watcher/common/log"
)

type FeedRoundHandler struct {
	listener       net.Listener
	connectedSocks []net.Conn
	mutex          sync.Mutex
}

// Client로부터 connect 요청을 대기할 port 바인딩 값이 필요함
func NewFeedRoundHandler(bindingAddress string) *FeedRoundHandler {
	listener, err := net.Listen("tcp", bindingAddress)
	if err != nil {
		LOG.Fatal("socket listening error : ", err.Error())
		os.Exit(-1)
	}
	return &FeedRoundHandler{listener: listener, mutex: sync.Mutex{}}
}

// conn객체를 이용해서 streamBuffer에 저장된 바이트 데이터를 보냅니다.
func (s *FeedRoundHandler) Send(jsonBytes json.RawMessage) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for _, sock := range s.connectedSocks {
		encoder := json.NewEncoder(sock)
		err := encoder.Encode(jsonBytes)
		if err != nil {
			LOG.Error("Encode fail.:", err.Error())
			return err
		}
		return nil
	}

	s.mutex.Unlock()
	return nil
}

// Client로부터 connect을 대기함
func (s *FeedRoundHandler) WaitFor() {
	go func() {
		conn, err := s.listener.Accept()
		if err != nil {
			LOG.Fatal("socket Accept error : ", err.Error())
			conn.Close()
			os.Exit(-1)
		}
		s.mutex.Lock()
		defer s.mutex.Unlock()
		s.connectedSocks = append(s.connectedSocks, conn)
		s.mutex.Unlock()
	}()
}
