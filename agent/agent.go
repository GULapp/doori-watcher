package main

import (
	"encoding/json"
	"fmt"
	"go_monitoring/agent/system"
	"log"
	"os"
)

const (
	KLOGFLAG = log.Lshortfile | log.Lmicroseconds
)

var (
	logFile *os.File
	logger  *log.Logger
	//logBuffer bytes.Buffer
	//logger = common.New(&logBuffer, "Agent : " , KLOGFLAG)
)

func init() {
	var err error
	logFile, err = os.OpenFile("/tmp/common", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Errorf("failed to open a logfile")
		os.Exit(-1)
	}
	logger = log.New(logFile, "Agent : ", KLOGFLAG)
}

func main() {
	var collection []Gather
	collection = append(collection, &system.Cpu{
		Totalusermode:   0,
		Totalsystemmode: 0,
		Totalnice:       0,
		Totalidle:       0,
		Totalwait:       0,
		Totalirq:        0,
		Totalsrq:        0,
		Cores:           nil,
	})
	for _, g := range collection {
		g.Gathering()

		b, err := json.MarshalIndent(g, "", "\t")
		if err != nil {
			log.Println("json Marshaling error")
			return
		}
		logger.Print(string(b))
		n, err := logFile.Write([]byte("leejaeseong\n"))
		if err != nil {
			fmt.Printf("Write err %d\n", n)
			fmt.Printf("Write err %s\n", err.Error())
		}
	}
}
