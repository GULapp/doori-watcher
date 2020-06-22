package common

import (
	"log"
	"net"
	"os"
)

type ReceiverHandler func([]byte)

type DrReceiver struct {
	bytesBuffer chan []byte
	conn        *net.IPConn
	addr        *net.IPAddr
	err         error
}

func NewDrReceiver() *DrReceiver {
	return &DrReceiver{bytesBuffer: make(chan []byte)}
}

func (s *DrReceiver) PullDataOnAsync(handler ReceiverHandler) {
	handler(<-s.bytesBuffer)
}

func (s *DrReceiver) WaitForConnection(network string, address string) {

	listener, err := net.Listen(network, address)
	if err != nil {
		log.Panicln(err.Error())
		os.Exit(-1)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Panicln(err.Error())
			conn.Close()
			os.Exit(-1)
		}
		go func(c net.Conn) {
			var bytes []byte
			readLen, err := c.Read(bytes)
			if err != nil {
				log.Panicln(err.Error())
			}
			if readLen == 0 {
				c.Close()
			}
			s.bytesBuffer <- bytes
		}(conn)
	}
}
