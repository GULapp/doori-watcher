package main

import (
	"go_monitoring/agent/system"
	"go_monitoring/collector"
	DRLOG "go_monitoring/common/log"
	"os"
	"time"
)

var (
	logging *DRLOG.DrLog
)

func init() {
	logging = DRLOG.NewDrLog("./agent.log", 0744)
}

func main() {
	logging.Info("START")

	//모니터링 대상 객체등록
	//Cpu, Ram 정보를 가지고 오는 객체등록
	//gather.go 에 정의된 인터페이스 정의 되어 있어햠
	//gathering 데이터를 수집하고, 결과를 []byte 리턴하는 함수
	//Done []byte를 입력받아 그 다음을 처리하는 함수
	var collection []collector.Gather
	collection = append(collection, &system.Cpu{})

	for{
		for _, g := range collection {
			stream, err := g.Gathering() //데이터수집
			if err != nil {
				logging.Fatal("gathering failed. %s", err.Error())
				os.Exit(-1)
			}
			g.Done(stream) //수집된 데이터를 처리.
		}
		time.Sleep(100*time.Millisecond)
	}
	os.Exit(0)
}
