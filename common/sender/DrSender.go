package sender

import (
	"log"
	"net"
)

const (
	kTcp = "tcp"
	kUdp = "udp"
)

type DrSender struct {
	bytesBuffer chan []byte
	conn        *net.IPConn
	addr        *net.IPAddr
	err         error
}

func NewDrSender() *DrSender {
	return &DrSender{bytesBuffer: make(chan []byte)}
}

func (s *DrSender) PushDataOnAsync(data []byte) {
	s.bytesBuffer <- data
}

func (s DrSender) Send() {
	sendingBuffer := <-s.bytesBuffer
	sendLen, err := s.conn.Write(sendingBuffer)
	if err != nil {
		log.Panicln("failed to IPConn.Write() ")
	}
	if sendLen != len(sendingBuffer) {
		log.Fatalln("wrong size")
	}
	log.Println("Send Data : ", sendingBuffer)
}

func (s *DrSender) SetIPAddr(network string, address string) {
	s.addr, s.err = net.ResolveIPAddr(network, address)
	if s.err != nil {
		log.Fatal("failed to set IP, addr")
	}
	log.Println("Address : ", s.addr.String())
}

func (s *DrSender) Connect() {
	s.conn, s.err = net.DialIP(s.addr.Network(), nil, s.addr)
	if s.err != nil {
		log.Panicln("failed to call DialIP ")
	}
}
