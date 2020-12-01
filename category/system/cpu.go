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
	var oldCpu Cpu
	oldCores := make([]Core, len(c.Cores))

	copy(oldCores, c.Cores)

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
				// 총 cpu 사용률 계산
				totalJiffes := Usermode + Systemmode + Nice + Idle
				totalJiffes -= oldCpu.Usermode + oldCpu.Systemmode + oldCpu.Nice + oldCpu.Idle
				totalIdleJiffes := Idle - oldCpu.Idle
				c.CpuUsage = int(100 * (1.0 - float32(totalIdleJiffes)/float32(totalJiffes)))
				oldCpu = Cpu{CpuUsage: c.CpuUsage, Usermode: Usermode, Systemmode: Systemmode, Nice: Nice, Idle: Idle}
			} else {
				if len(oldCores) < i {
					oldCores = append(oldCores, Core{})
				}
				//각 cpu core당 사용률 계산
				oldCores[i-1].Corename = cpuName
				totalJiffes := Usermode + Systemmode + Nice + Idle
				totalJiffes -= oldCores[i-1].Usermode + oldCores[i-1].Systemmode + oldCores[i-1].Nice + oldCores[i-1].Idle
				totalIdleJiffes := Idle - oldCores[i-1].Idle
				oldCores[i-1].Coreusage = int(100 * (1.0 - float32(totalIdleJiffes)/float32(totalJiffes)))
				oldCores[i-1] = Core{Corename: cpuName, Coreusage: oldCores[i-1].Coreusage, Usermode: Usermode, Systemmode: Systemmode, Nice: Nice, Idle: Idle}
			}
		}
	}
	c.Cores = nil
	c.Cores = make([]Core,len(oldCores))
	copy(c.Cores, oldCores)
	oldCores = nil
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
