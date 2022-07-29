package variables

import (
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/host"
	"os"
	"time"
)

type PeerInfo struct {
	Id       string    `json:"node_id"`
	UUID     string    `json:"uuid"`
	Ip       string    `json:"ip"`
	Hostname string    `json:"host,omitempty"`
	OS       string    `json:"os"`
	Platform string    `json:"platform"`
	Version  string    `json:"version"`
	Time     time.Time `json:"time"`
}

type Cpu struct {
	Id       string    `json:"node_id"`
	UUID     string    `json:"uuid"`
	Ip       string    `json:"ip"`
	Hostname string    `json:"host,omitempty"`
	Model    string    `json:"model"`
	Usage    int       `json:"usage, omitempty"`
	Time     time.Time `json:"time, omitempty"`
}

type Ram struct {
	Id       string    `json:"node_id"`
	UUID     string    `json:"uuid"`
	Ip       string    `json:"ip"`
	Hostname string    `json:"host,omitempty"`
	Usage    int       `json:"usage, omitempty"`
	Time     time.Time `json:"time, omitempty"`
}

type PingStatus struct {
	UUID   string    `json:"uuid"`
	Source string    `json:"source"`
	Target string    `json:"target"`
	Alive  bool      `json:"alive"`
	RTT    int64     `json:"rtt"`
	Time   time.Time `json:"time"`
}

// NewPeerInfo create method
func NewPeerInfo(ip, nodeId string) *PeerInfo {

	var platform, version, os = platformInformation()
	var host = hostname()

	return &PeerInfo{
		Id:       nodeId,
		Ip:       ip,
		Hostname: host,
		Platform: platform,
		Version:  version,
		OS:       os,
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

	model := cpuModel()
	host := hostname()

	return &Cpu{
		Id:       nodeId,
		Ip:       ip,
		Hostname: host,
		Model:    model,
	}
}

func hostname() string {
	hostName, err := os.Hostname()
	if err != nil {
		return ""
	}
	return hostName
}

//Return different string as platform, version und os
func platformInformation() (string, string, string) {
	hostStat, err := host.Info()
	if err != nil {
		return "", "", ""
	}
	platform := hostStat.Platform
	version := hostStat.PlatformVersion
	os := hostStat.OS

	return platform, version, os
}

//return Cpu Model
func cpuModel() string {
	cpuStat, err := cpu.Info()
	if err != nil {
		return ""
	}
	return cpuStat[0].ModelName
}
