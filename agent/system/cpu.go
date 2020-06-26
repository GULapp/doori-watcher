package system

import (
	"encoding/json"
	LOG "go_monitoring/common/log"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type Cpu struct {
	Totalusermode   uint32 `json:"totalusermode"`
	Totalsystemmode uint32 `json:"totalsystemmode"`
	Totalnice       uint32 `json:"-"`
	Totalidle       uint32 `json:"-"`
	Totalwait       uint32 `json:"-"`
	Totalirq        uint32 `json:"-"`
	Totalsrq        uint32 `json:"-"`
	Cores           []cpuCore
}

type cpuCore struct {
	Corenum    uint32 `json:corenum`
	Usermode   uint32 `json:"usermode"`
	Systemmode uint32 `json:"systemmode"`
	Nice       uint32 `json:"-"`
	Idle       uint32 `json:"-"`
	Wait       uint32 `json:"-"`
	Irq        uint32 `json:"-"`
	Srq        uint32 `json:"-"`
}

func (c *Cpu) PrettyPrint() {
	LOG.Debug("Total System Cpu Usage %d", c.Totalsystemmode)
	LOG.Debug("Total User Cpu Usage %d", c.Totalusermode)
	LOG.Debug("Total nice %d", c.Totalnice)
	LOG.Debug("Total idle %d", c.Totalidle)
	LOG.Debug("Total wait %d", c.Totalwait)
	LOG.Debug("Total IRQ %d", c.Totalirq)
	LOG.Debug("Total SRQ %d", c.Totalsrq)

	for i, coreMembers := range c.Cores {
		LOG.Debug("core[%d] System Usage %d", i, coreMembers.Systemmode)
		LOG.Debug("core[%d] User Usage %d, below", i, coreMembers.Usermode)
		LOG.Debug("nice %d", coreMembers.Nice)
		LOG.Debug("idle %d", coreMembers.Idle)
		LOG.Debug("wait %d", coreMembers.Wait)
		LOG.Debug("IRQ %d", coreMembers.Irq)
		LOG.Debug("SRQ %d", coreMembers.Srq)
	}
}
func (c *Cpu) clear() {
	c.Cores = nil
}

func (c *Cpu) Gathering() []byte {
	c.clear()

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
		fields := strings.Fields(line)
		if strings.Contains(fields[0], "cpu") {
			if i == 0 {
				Totalusermode, _ := strconv.ParseUint(fields[1], 10, 32)
				Totalsystemmode, _ := strconv.ParseUint(fields[2], 10, 32)
				Totalnice, _ := strconv.ParseUint(fields[3], 10, 32)
				Totalidle, _ := strconv.ParseUint(fields[4], 10, 32)
				Totalwait, _ := strconv.ParseUint(fields[5], 10, 32)
				Totalirq, _ := strconv.ParseUint(fields[6], 10, 32)
				Totalsrq, _ := strconv.ParseUint(fields[7], 10, 32)

				c.Totalusermode = uint32(Totalusermode)
				c.Totalsystemmode = uint32(Totalsystemmode)
				c.Totalnice = uint32(Totalnice)
				c.Totalidle = uint32(Totalidle)
				c.Totalwait = uint32(Totalwait)
				c.Totalirq = uint32(Totalirq)
				c.Totalsrq = uint32(Totalsrq)
			} else {
				var cpuCore cpuCore
				Usermode, _ := strconv.ParseUint(fields[1], 10, 32)
				Systemmode, _ := strconv.ParseUint(fields[2], 10, 32)
				Nice, _ := strconv.ParseUint(fields[3], 10, 32)
				Idle, _ := strconv.ParseUint(fields[4], 10, 32)
				Wait, _ := strconv.ParseUint(fields[5], 10, 32)
				Irq, _ := strconv.ParseUint(fields[6], 10, 32)
				Srq, _ := strconv.ParseUint(fields[7], 10, 32)

				cpuCore.Corenum = uint32(i)
				cpuCore.Usermode = uint32(Usermode)
				cpuCore.Systemmode = uint32(Systemmode)
				cpuCore.Nice = uint32(Nice)
				cpuCore.Idle = uint32(Idle)
				cpuCore.Wait = uint32(Wait)
				cpuCore.Irq = uint32(Irq)
				cpuCore.Srq = uint32(Srq)
				c.Cores = append(c.Cores, cpuCore)
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
