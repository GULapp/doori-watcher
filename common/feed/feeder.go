package feed

import (
	"bufio"
	"encoding/binary"
	"io"
	"net"
	"os"
	LOG "watcher/common/log"
)

type ProcHandler func(data <-chan []byte)

type Feeder struct {
	DataEventChan chan []byte
	procData      ProcHandler
	conn          *net.Conn
	err           error
}

/*Feeder로 들어온 데이터를 처리할 함수 type ProcHandler func(chan []byte) 형의 인수로 받아야 함*/
func NewFeeder(handler ProcHandler) *Feeder {
	return &Feeder{DataEventChan: make(chan []byte), procData: handler}
}

/*chan 넘겨서, 해당 ProcHandler에게 처리를 위임함*/
func (f *Feeder) BringupFeeder() {
	go f.procData(f.DataEventChan)
}

func (f *Feeder) WaitFor(commType string, address string) {
	switch commType {
	case "tcp":
		f.waitForTcp(address)
		break
	default:
		break
	}
}

func (f *Feeder) waitForTcp(address string) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		LOG.Fatal("socket listening error : ", err.Error())
		os.Exit(-1)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			LOG.Fatal("socket Accept error : ", err.Error())
			conn.Close()
			os.Exit(-1)
		}

		go func(conn net.Conn) {
			reader:=bufio.NewReader(conn)
			for {
				//들어온 데이터의 총길이가 다 도착할까지 blocking처리
				dataBufferLen := make([]byte,2)
				n, err :=io.ReadFull(reader,dataBufferLen)
				if err != nil || n != 2 {
					LOG.Error("Read() error %s", err.Error())
					conn.Close()
					return
				}
				//나머지 데이터가 들어올때까지 대기
				bodyLength := binary.LittleEndian.Uint16(dataBufferLen)
				dataBufferBody := make([]byte,bodyLength)
				n, err = io.ReadFull(reader,dataBufferBody)
				if err != nil || uint16(n) != bodyLength {
					LOG.Error("Read() error %s", err.Error())
					conn.Close()
					return
				}
				f.DataEventChan <- dataBufferBody /*chan 데이터 보내기, 딱 받은 사이즈만큼만, [:n]*/
			}
		}(conn)
	}
}
