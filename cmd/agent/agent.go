package main

import (
	"encoding/json"
	"net"
	"os"
	"time"
	"watcher/category/system"
	"watcher/collector"
	"watcher/common"
	"watcher/common/config"
	"watcher/common/feed"
	LOG "watcher/common/log"
)

var (
	tomlConfig config.Config
	gFeedHandler *feed.FeedHandler
	gConn        net.Conn
)

func init() {
	var err error
	tomlConfig, err = config.InitConfig("./config.toml")
	if err != nil {
		LOG.Fatal("InitConfig() error, process is terminated")
		os.Exit(-1)
	}

	//LOG initialize.
	LOG.Init("/tmp/agent.log", LOG.TRACE, 0744)
}

func main() {
	LOG.Info("System Monitoring Agent START")


	if err := connectFeeder(); err != nil {
		LOG.Fatal("failed to connect Feeder %s", err.Error())
		os.Exit(-1)
	}

	bytesQueue := make(chan json.RawMessage)

	//모니터링 대상 객체등록
	//Cpu, Ram 정보를 가지고 오는 객체등록
	//gather.go 에 정의된 인터페이스 정의 되어 있어햠
	//gathering 데이터를 수집하고, 결과를 []byte 리턴하는 함수
	//Done []byte를 입력받아 그 다음을 처리하는 함수
	var collection []collector.Gather
	collection = append(collection, &system.Cpu{})

	/* monitoring server */
	go func() {
		for {
			buffer := <-bytesQueue
			if err := gFeedHandler.Send(gConn, buffer); err != nil {
				LOG.Fatal("connecting to server error : ", err.Error())
				gConn.Close()
				os.Exit(-1)
			}
		}
	}()

	protocol := common.Protocol{}
	protocol.Init(tomlConfig.Site,tomlConfig.Domain,"local","127.0.0.1","system")

	for {
		for _, g := range collection {
			outputJson := g.Gathering() //데이터수집
			g.PrettyPrint()
			protocol.Set("Cpu", outputJson)
			marshalingBytes, _ := protocol.Marshaling()
			bytesQueue <- marshalingBytes
			g.Done(outputJson) //수집된 데이터를 처리.
		}
		time.Sleep(1000 * time.Millisecond) /*1 second sleep*/
	}

	LOG.Info("Agent Process is terminated.")
	os.Exit(0)
}

func connectFeeder() error {
	var err error
	gFeedHandler = feed.NewFeedHandler()
	//to Feeder
	gConn, err = gFeedHandler.Connect(tomlConfig.Agent.To.Ip+":"+tomlConfig.Agent.To.Port)
	if err != nil {
		LOG.Fatal("failed to Call FeedHandler : %s", err.Error())
		return err
	}
	return nil
}
