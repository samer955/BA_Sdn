package variables

import (
	"os"
	"time"
)

type PeerInfo struct {
	Id       string    `json:"node_id"`
	Ip       string    `json:"ip"`
	Hostname string    `json:"host,omitempty"`
	Time     time.Time `json:"time"`
}

type Cpu struct {
	Id        string    `json:"node_id"`
	Ip        string    `json:"ip"`
	Processes []Process `json:"processes"`
	Hostname  string    `json:"host,omitempty"`
	Usage     int       `json:"usage, omitempty"`
	Time      time.Time `json:"time, omitempty"`
}

type Ram struct {
	Id       string    `json:"node_id"`
	Ip       string    `json:"ip"`
	Hostname string    `json:"host,omitempty"`
	Usage    int       `json:"usage, omitempty"`
	Time     time.Time `json:"time, omitempty"`
}

type Process struct {
	Name        string  `json:"name"`
	Cpu_percent float64 `json:"cpu_percent"`
}

type PingStatus struct {
	Source_node string    `json:"from"`
	Target_node string    `json:"target"`
	Alive       bool      `json:"alive"`
	RTT         int64     `json:"rtt"`
	Time        time.Time `json:"time"`
}

// NewPeerInfo create method
func NewPeerInfo(ip, nodeId string) *PeerInfo {
	return &PeerInfo{
		Id:       nodeId,
		Ip:       ip,
		Hostname: hostname(),
	}
}

// Ram create method
func NewRam(ip, nodeId string) *Ram {
	return &Ram{
		Id:       nodeId,
		Ip:       ip,
		Hostname: hostname(),
	}
}

// CPU create method
func NewCpu(ip, nodeId string) *Cpu {
	return &Cpu{
		Id:       nodeId,
		Ip:       ip,
		Hostname: hostname(),
	}
}

func hostname() string {
	hostName, err := os.Hostname()
	if err != nil {
		return ""
	}
	return hostName

}
