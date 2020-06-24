package system

import (
	"encoding/json"
	DRLOG "go_monitoring/common/log"
	DRSENDER "go_monitoring/common/sender"
	"io/ioutil"
	"net"
	"os"
	"strconv"
	"strings"
)

var (
	logging *DRLOG.DrLog
	socket  *DRSENDER.DrSender
	fd      net.Conn
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
	Usermode   uint32 `json:"usermode"`
	Systemmode uint32 `json:"systemmode"`
	Nice       uint32 `json:"-"`
	Idle       uint32 `json:"-"`
	Wait       uint32 `json:"-"`
	Irq        uint32 `json:"-"`
	Srq        uint32 `json:"-"`
}

func init() {
	logging = DRLOG.NewDrLog("./agent.log", 0744)
	socket = DRSENDER.NewDrSender()
	var err error
	if fd, err = socket.Connect("localhost:12345"); err != nil {
		logging.Info("failed to connect to SERVER")
		os.Exit(-1)
	}
}

func (c *Cpu) PrettyPrint() {
	logging.Debug("%d", c.Totalsystemmode)
	logging.Debug("%d", c.Totalusermode)
	logging.Debug("%d", c.Totalusermode)
	logging.Debug("%d", c.Totalsystemmode)
	logging.Debug("%d", c.Totalnice)
	logging.Debug("%d", c.Totalidle)
	logging.Debug("%d", c.Totalwait)
	logging.Debug("%d", c.Totalirq)
	logging.Debug("%d", c.Totalsrq)

	for i, coreMembers := range c.Cores {
		logging.Debug("cpu core[%d] %d", i, coreMembers.Usermode)
		logging.Debug("%d", coreMembers.Systemmode)
		logging.Debug("%d", coreMembers.Nice)
		logging.Debug("%d", coreMembers.Idle)
		logging.Debug("%d", coreMembers.Wait)
		logging.Debug("%d", coreMembers.Irq)
		logging.Debug("%d", coreMembers.Srq)
	}
}

func (c *Cpu) Gathering() []byte {
	contents, err:= ioutil.ReadFile("/proc/stat")
	if err != nil {
		logging.Fatal("cant read /proc/stat : %s", err.Error())
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
	jsonBytes, err := json.MarshalIndent(c, "", "")
	if err != nil {
		logging.Fatal("Marshaling error : %s", err.Error())
		os.Exit(-1)
	}
	return jsonBytes
}

/* Socket을 이용하여(already Established) 데이터 송신 */
func (c *Cpu) Done(buffer []byte) error {
	logging.Debug("finished sending :", buffer)
	if err := socket.Send(fd, buffer); err != nil {
		logging.Fatal("Cpu, Send() error : %s", err.Error())
		return err
	}
	return nil
}
