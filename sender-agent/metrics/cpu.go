package metrics

import (
	"fmt"
	"github.com/shirou/gopsutil/v3/cpu"
	"os"
	"time"
)

type Cpu struct {
	Id       string    `json:"node_id"`
	UUID     string    `json:"uuid"`
	Ip       string    `json:"ip"`
	Hostname string    `json:"host,omitempty"`
	Model    string    `json:"model"`
	Usage    int       `json:"usage, omitempty"`
	Time     time.Time `json:"time, omitempty"`
}

// CPU create method
func NewCpu(ip, nodeId string) *Cpu {

	hostname, err := os.Hostname()
	if err != nil {
		hostname = ""
	}

	model := cpuModel()

	return &Cpu{
		Id:       nodeId,
		Ip:       ip,
		Hostname: hostname,
		Model:    model,
	}
}

//return Cpu Model
func cpuModel() string {
	cpuStat, err := cpu.Info()
	if err != nil {
		return "Not available"
	}
	return cpuStat[0].ModelName
}

//Get the actual CPU Percentage from the system
func (c *Cpu) UpdateCpuPercentage() {
	c.Time = time.Now()
	cpuUsage, err := cpu.Percent(0, false)
	if err != nil {
		fmt.Println("Unable to get Cpu Percentage")
		c.Usage = 00
		return
	}
	c.Usage = int(cpuUsage[0])
}
