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

func NewFeedRoundHandler(bindingAddress string) *FeedRoundHandler {
	listener, err := net.Listen("tcp", bindingAddress)
	if err != nil {
		LOG.Fatal("socket listening error : ", err.Error())
		os.Exit(-1)
	}
	return &FeedRoundHandler{listener: listener, mutex: sync.Mutex{}}
}

// conn객체를 이용해서 streamBuffer에 저장된 바이트 데이터를 보냅니다.
func (s *FeedRoundHandler) Send(jsonBytes DataContainer) error {
	s.mutex.Lock()

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

// address에 tcp 연결합니다.
func (s *FeedRoundHandler) WaitFor() {
	go func() {
		conn, err := s.listener.Accept()
		if err != nil {
			LOG.Fatal("socket Accept error : ", err.Error())
			conn.Close()
			os.Exit(-1)
		}
		s.mutex.Lock()
		s.connectedSocks = append(s.connectedSocks, conn)
		s.mutex.Unlock()
	}()
}
