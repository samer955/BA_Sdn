package metrics

import (
	"fmt"
	"github.com/shirou/gopsutil/v3/mem"
	"os"
	"time"
)

type Ram struct {
	Id       string    `json:"node_id"`
	UUID     string    `json:"uuid"`
	Ip       string    `json:"ip"`
	Hostname string    `json:"host,omitempty"`
	Usage    int       `json:"usage, omitempty"`
	Time     time.Time `json:"time, omitempty"`
}

// Ram create method
func NewRam(ip, nodeId string) *Ram {

	host, _ := os.Hostname()

	return &Ram{
		Id:       nodeId,
		Ip:       ip,
		Hostname: host,
	}
}

//Get the actual RAM Percentage from the system
func (ram *Ram) UpdateRamPercentage() {
	ram.Time = time.Now()
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		fmt.Println("Unable to get Memory Info")
		ram.Usage = 00
		return
	}
	ram.Usage = int(vmStat.UsedPercent)
}
