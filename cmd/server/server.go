package main

import (
	"bufio"
	DRLOG "go_monitoring/common/log"
	"net"
)

var (
	logging *DRLOG.DrLog
)

func init() {
	logging = DRLOG.NewDrLog("./server.log", 0744)
}

func main() {
	logging.Info("Monitoring SERVER START")

	listener, err := net.Listen("tcp", "localhost:12345")
	if err != nil {
		logging.Fatal("Server Monitoring socket, Listening error %s", err.Error())
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			logging.Error("Accept error, %s", err.Error())
		}
		go func(conn net.Conn) {
			buffers := make([]byte, 4096)
			for {
				n, err := bufio.NewReader(conn).Read(buffers)
				if err != nil {
					logging.Error("Read() error %s", err.Error())
				}
				logging.Info("Read Len:%d", n)
				if n > 0 {
					logging.Info(string(buffers))
				}
			}
		}(conn)
	}
}
