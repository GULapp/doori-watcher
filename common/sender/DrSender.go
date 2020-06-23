package sender

import (
	"net"
	DRLOG "go_monitoring/common/log"
)

var (
	logging *DRLOG.DrLog
)

const (
	TCP = "tcp"
	UDP = "udp"
)

type DrSender struct {
	conn *net.Dialer
}

func init() {
	logging = DRLOG.NewDrLog("./agent.log", 0744)
}

func NewDrSender() *DrSender {
	return &DrSender{}
}

func (s DrSender) Send(streamBuffer []byte) error {
	sendLen, err := s.conn.Write(streamBuffer)
	if err != nil {
		logging.Fatal("failed to Write(), %s", err.Error())
		return err
	}
	if sendLen != len(streamBuffer) {
		logging.Fatal("wrong size")
		return err
	}
	logging.Debug("Send Data : ", streamBuffer)
	return nil
}

func (s *DrSender) Connect(network string, address string) error{
	//var err error
	//s.conn, err = s.conn.Dial(network, address)
	//if err != nil {
	//	logging.Fatal("failed to call DialIP ")
	//	return err
	//}
	return nil
}
