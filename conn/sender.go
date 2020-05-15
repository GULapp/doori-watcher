package conn

import (
	"log"
	"net"
)

const (
	TCP = "tcp"
	UDP = "udp"
)

type Sender struct{
	bytesBuffer chan []byte
	conn net.IPConn
	addr net.IPAddr
	err error
}

func NewSender() *Sender{
	return &Sender{bytesBuffer: make(chan []byte)}
}

func (s *Sender) PushDataOnAsync(data []byte) {
	s.bytesBuffer <- data
}

func (s Sender) Send() {
	sendingBuffer := <- s.bytesBuffer
	sendLen, err := s.conn.Write( sendingBuffer )
	if err != nil {
		log.Panicln("failed to IPConn.Write() ")
	}
	if sendLen != len(sendingBuffer) {
		log.Fatalln("wrong size")
	}
	log.Println("Send Data : ", sendingBuffer)
}

func (s *Sender) SetIPAddr(network string , address string) {
	s.addr, s.err = net.ResolveIPAddr(network, address)
	if s.err != nil {
		log.Fatal("failed to set IP, addr")
	}
	log.Println("Address : ", s.addr.String() )
}

func (s *Sender) Connect() {
	s.conn, s.err = net.DialIP( s.addr.Network(), nil, &s.addr )
	if s.err != nil {
		log.Panicln("failed to call DialIP ")
	}
}
