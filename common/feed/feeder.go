package feed

import (
	"bufio"
	"math"
	"net"
	"os"
	LOG "watcher/common/log"
)

type ProcHandler func(data <-chan []byte)

type Feeder struct {
	bytesBuffer chan []byte
	procData 	ProcHandler
	conn        *net.Conn
	err         error
}

/*Feeder로 들어온 데이터를 처리할 함수 type ProcHandler func(chan []byte) 형의 인수로 받아야 함*/
func NewFeeder(handler ProcHandler) *Feeder {
	return &Feeder{bytesBuffer: make(chan []byte), procData:handler}
}

/*chan 넘겨서, 해당 ProcHandler에게 처리를 위임함*/
func (f *Feeder) BringupFeeder() {
	go f.procData(f.bytesBuffer)
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
			buffers := make([]byte, math.MaxInt16)
			for {
				n, err := bufio.NewReader(conn).Read(buffers)
				if err != nil {
					LOG.Error("Read() error %s", err.Error())
					conn.Close()
					return
				}
				f.bytesBuffer<-buffers[:n] /*chan 데이터 보내기, 딱 받은 사이즈만큼만, [:n]*/
			}
		}(conn)
	}
}
