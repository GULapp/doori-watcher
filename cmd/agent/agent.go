package main

import (
	"net"
	"os"
	"time"
	"watcher/collector"
	"watcher/common/feed"
	LOG "watcher/common/log"
	"watcher/objects/system"
)
var(
	gFeedHandler *feed.FeedHandler
	gConn        net.Conn
)

func init() {
	//LOG initialize.
	LOG.Init("/tmp/agent.log", LOG.ERROR,0744)
}

func connectFeeder() error{
	var err error
	gFeedHandler = feed.NewFeedHandler()
	//to Feeder
	gConn, err = gFeedHandler.Connect("localhost:12345")
	if err != nil {
		LOG.Fatal("failed to Call FeedHandler : %s", err.Error())
		return err
	}
	return nil
}

func main() {
	LOG.Info("System Monitoring Agent START")

	if err := connectFeeder(); err != nil {
		LOG.Fatal("failed to connect Feeder %s",err.Error())
		os.Exit(-1)
	}

	bytesQueue := make(chan feed.DataContainer)

	//모니터링 대상 객체등록
	//Cpu, Ram 정보를 가지고 오는 객체등록
	//gather.go 에 정의된 인터페이스 정의 되어 있어햠
	//gathering 데이터를 수집하고, 결과를 []byte 리턴하는 함수
	//Done []byte를 입력받아 그 다음을 처리하는 함수
	var collection []collector.Gather
	collection = append(collection, &system.Cpu{})

	/* monitoring server */
	go func() {
		for{
			buffer :=<-bytesQueue
			if err:= gFeedHandler.Send(gConn, buffer); err!= nil {
				LOG.Fatal("connecting to server error : ", err.Error())
				gConn.Close()
				os.Exit(-1)
			}
		}
	}()

	for {
		for _, g := range collection {
			outputJson := g.Gathering() //데이터수집
			g.PrettyPrint()
			bytesQueue <- feed.DataContainer{"Cpu", outputJson}
			g.Done(outputJson) //수집된 데이터를 처리.
		}
		time.Sleep(1000*time.Millisecond) /*1 second sleep*/
	}

	LOG.Info("Agent Process is terminated.")
	os.Exit(0)
}
