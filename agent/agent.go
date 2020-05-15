package main

import (
	"agent/system"
	"bytes"
	"encoding/json"
	"log"
)

const (
	KLOGFLAG = log.Lshortfile | log.Lmicroseconds
	//KLOGFLAG = log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile
)


var (
	logBuffer bytes.Buffer
	logger = log.New(&logBuffer, "Agent : " , KLOGFLAG)
)

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

		b,err := json.MarshalIndent( g, "", "\t" )
		if err != nil {
			log.Println("json Marshaling error")
			return
		}
		logger.Print(string(b))
	}
}
