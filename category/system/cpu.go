package system

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"watcher/common"
	LOG "watcher/common/log"
)

// json:"-" 는, Marshaling할때, 무시됨. 즉 추가되지 않음.
// json Marshaling를 위해서는 멤버는 exported 형태(노출되어야) 대문자여야 한다.
type Cpu struct {
	CpuUsage   int `json:"total usage"`
	Usermode   int `json:"-"`
	Systemmode int `json:"-"`
	Nice       int `json:"-"`
	Idle       int `json:"-"`
	Wait       int `json:"-"`
	Irq        int `json:"-"`
	Srq        int `json:"-"`
	Cores      []Core
}

type Core struct {
	Corename   string `json:"core name"`
	Coreusage  int    `json:"usage"`
	Usermode   int    `json:"-"`
	Systemmode int    `json:"-"`
	Nice       int    `json:"-"`
	Idle       int    `json:"-"`
	Wait       int    `json:"-"`
	Irq        int    `json:"-"`
	Srq        int    `json:"-"`
}

func (c *Cpu) PrettyPrint() {
	common.PrintAsStructForJson(c)
}

func (c *Cpu) Gathering() []byte {
	contents, err := ioutil.ReadFile("/proc/stat") /*Linux cpu 정보*/
	if err != nil {
		LOG.Fatal("cant read /proc/stat : %s", err.Error())
		os.Exit(-1)
	}
	lines := strings.Split(string(contents), "\n")
	for i, line := range lines {
		if lines[i] == "" {
			break
		}
		io := strings.NewReader(line)
		var cpuName string
		var Usermode, Systemmode, Nice, Idle, Wait, Irq, Srq int
		var Coreusage int
		fmt.Fscanf(io, "%s %d %d %d %d %d %d %d",
			&cpuName,
			&Usermode,
			&Systemmode,
			&Nice,
			&Idle,
			&Wait,
			&Irq,
			&Srq)

		if strings.Contains(cpuName, "cpu") {
			if i == 0 {
				// 총 cpu 사용률 계산, c. prefix값은 oldcpu 값
				totalJiffes := Usermode + Systemmode + Nice + Idle
				totalJiffes -= c.Usermode + c.Systemmode + c.Nice + c.Idle
				totalIdleJiffes := Idle - c.Idle
				c.CpuUsage = int(100 * (1.0 - float32(totalIdleJiffes)/float32(totalJiffes)))
				c.Usermode  = Usermode
				c.Systemmode= Systemmode
				c.Nice      = Nice
				c.Idle      = Idle
				c.Wait      = Wait
				c.Irq       = Irq
				c.Srq       = Srq
			} else {
				if len(c.Cores) < i {
					c.Cores = append(c.Cores, Core{})
				}
				//각 cpu core당 사용률 계산
				c.Cores[i-1].Corename = cpuName
				totalJiffes := Usermode + Systemmode + Nice + Idle
				totalJiffes -= c.Cores[i-1].Usermode + c.Cores[i-1].Systemmode + c.Cores[i-1].Nice + c.Cores[i-1].Idle
				totalIdleJiffes := Idle - c.Cores[i-1].Idle
				Coreusage = int(100 * (1.0 - float32(totalIdleJiffes)/float32(totalJiffes)))
				c.Cores[i-1] = Core{Corename: cpuName, Coreusage: Coreusage, Usermode: Usermode, Systemmode: Systemmode, Nice: Nice, Idle: Idle}
			}
		}
	}
	return c.serialize()
}

func (c *Cpu) serialize() []byte {
	jsonBytes, err := json.Marshal(c)
	if err != nil {
		LOG.Fatal("Marshaling error : %s", err.Error())
		os.Exit(-1)
	}
	return jsonBytes
}

func (c *Cpu) Done(buffer []byte) error {
	LOG.Info("Cpu gathering, done %s", string(buffer))
	return nil
}
