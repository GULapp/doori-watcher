package system

import (
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

type cpuCore struct {
	Usermode   uint32 `json:"usermode"`
	Systemmode uint32 `json:"systemmode"`
	Nice       uint32 `json:"nice"`
	Idle       uint32 `json:"idle"`
	Wait       uint32 `json:"wait"`
	Irq        uint32 `json:"irq"`
	Srq        uint32 `json:"srq"`
}

type Cpu struct {
	Totalusermode   uint32 `json:"totalusermode"`
	Totalsystemmode uint32 `json:"totalsystemmode"`
	Totalnice       uint32 `json:"totalnice"`
	Totalidle       uint32 `json:"totalidle"`
	Totalwait       uint32 `json:"totalwait"`
	Totalirq        uint32 `json:"totalirq"`
	Totalsrq        uint32 `json:"totalsrq"`
	Cores           []cpuCore
}

func (c *Cpu) Display() {
	log.Printf("%d", c.Totalsystemmode)
	log.Printf("%d", c.Totalusermode)
	log.Printf("%d", c.Totalusermode)
	log.Printf("%d", c.Totalsystemmode)
	log.Printf("%d", c.Totalnice)
	log.Printf("%d", c.Totalidle)
	log.Printf("%d", c.Totalwait)
	log.Printf("%d", c.Totalirq)
	log.Printf("%d", c.Totalsrq)

	for i, coreMembers := range c.Cores {
		log.Printf("cpu core[%d] %d", i, coreMembers.Usermode)
		log.Printf("%d", coreMembers.Systemmode)
		log.Printf("%d", coreMembers.Nice)
		log.Printf("%d", coreMembers.Idle)
		log.Printf("%d", coreMembers.Wait)
		log.Printf("%d", coreMembers.Irq)
		log.Printf("%d", coreMembers.Srq)
	}
}

func (c *Cpu) Gathering() error {
	contents, err := ioutil.ReadFile("/proc/stat")
	if err != nil {
		return err
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
	return nil
}
