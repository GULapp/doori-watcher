package sender

import (
	"bufio"
	"net"
	DRLOG "go_monitoring/common/log"
	"time"
)

var (
	logging *DRLOG.DrLog
)

const (
	TCP = "tcp"
	UDP = "udp"
)

type DrSender struct {
}

func init() {
	logging = DRLOG.NewDrLog("./agent.log", 0744)
}

func NewDrSender() *DrSender {
	return &DrSender{}
}

func (s *DrSender) Send(conn net.Conn, streamBuffer []byte) error {
	writedLen,err :=bufio.NewWriter(conn).Write(streamBuffer)
	if writedLen < len(streamBuffer) {
		logging.Fatal("wrong size, %s", err.Error())
		return err
	}
	return nil
}

func (s *DrSender) Connect(address string) (net.Conn, error) {
	conn, err := net.DialTimeout(TCP, address, 10*time.Millisecond)
	if err != nil {
		logging.Fatal("failed to call DialIP %s", err.Error())
		return nil, err
	}
	return conn, nil
}
