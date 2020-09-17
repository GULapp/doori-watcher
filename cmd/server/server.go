package main

import (
	"encoding/json"
	"os"
	"watcher/common"
	"watcher/category/system"

	//"encoding/json"
	"watcher/common/feed"
	LOG "watcher/common/log"
	//"watcher/category/system"
)

func init() {
	//LOG initialize.
	LOG.Init("/tmp/server.log", LOG.TRACE, 0744)
}

func main() {
	LOG.Info("Monitoring SERVER START")

	forClients := feed.FeedRoundHandler{}

	/*데이터가 들어오면, procTcpData 함수한테 처리하도록 위임. 콜백 등록함.*/
	forAgents := feed.NewFeeder(procTcpData)
	/*데이터가 오기를 기다리는 상태(채널 수신이벤트 기다림), procTcpData 호출됨*/
	forAgents.BringupFeeder()
	/*데이터를 수신용, 통신 열기*/
	forAgents.WaitFor("tcp", "localhost:12345")
}

func procTcpData(dataQueue <-chan json.RawMessage) {
	for {
		var inputMessage common.Protocol

		messages := <-dataQueue /*채널에서 데이터가 수신*/
		/*Protocol 구조체로 역직렬화*/
		err := json.Unmarshal(messages, &inputMessage)
		if err != nil {
			LOG.Error("Unmarshal error : %s", err.Error())
		}
		LOG.Debug("Receive Type:%s", inputMessage.Body.Tr)

		TR := inputMessage.Body.Tr

		switch TR {
		case "Cpu":
			var cpuInfo system.Cpu
			err = json.Unmarshal(inputMessage.Body.Data, &cpuInfo)
			if err != nil {
				LOG.Error("Unmarshal error : %s", err.Error())
			}
			cpuInfo.PrettyPrint()
		default:
			LOG.Error("unknown TR")
			os.Exit(-1)
		}
	}
}
