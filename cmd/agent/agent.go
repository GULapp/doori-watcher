package main

import (
	"encoding/json"
	"go_monitoring/agent/system"
	"go_monitoring/collector"
	DRLOG "go_monitoring/common/log"
)

var (
	logging *DRLOG.DrLog
)

func init() {
	logging = DRLOG.NewDrLog("./agent.log", 0744)
}

func main() {
	logging.Info("START")
	var collection []collector.Gather
	collection = append(collection, &system.Cpu{} )
	for _, g := range collection {
		g.Gathering()

		b, err := json.MarshalIndent(g, "", "\t")
		if err != nil {
			logging.Fatal("json Marshaling error")
			return
		}
		logging.Debug(string(b))
	}
}
