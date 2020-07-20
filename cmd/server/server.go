package main

import (
	"encoding/json"
	"watcher/common/feed"
	LOG "watcher/common/log"
	"watcher/objects/system"
)

func init() {
	//LOG initialize.
	LOG.Init("/tmp/server.log", LOG.ERROR, 0744)
}

func main() {
	LOG.Info("Monitoring SERVER START")

	/*데이터가 들어오면, procTcpData 함수한테 처리하도록 위임. 등록함.*/
	feeder := feed.NewFeeder(procTcpData)

	/*채널에서 데이터가 수신대기 상태, 데이터 수신시, procTcpData 호출됨*/
	feeder.BringupFeeder()

	/*소켓을 열고, agent로부터 연결을 기다린다. 연결이 완료되고, 데이터 수신이 되면, 채널로 데이터보내기 */
	feeder.WaitFor("tcp", "localhost:12345")
}

func procTcpData(dataQeueue <-chan []byte) {
	for {
		messages := <-dataQeueue /*채널에서 데이터가 수신되면 다시 message 변수에 전달*/
		LOG.Debug(string(messages))

		var cpuinfo system.Cpu
		/*역직렬화*/
		err := json.Unmarshal(messages[:len(messages)], &cpuinfo)
		if err != nil {
			LOG.Error("Unmarshal error : %s", err.Error())
		}
		cpuinfo.PrettyPrint()
	}
}