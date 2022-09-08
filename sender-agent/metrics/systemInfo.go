package metrics

import (
	"bufio"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/shirou/gopsutil/v3/host"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

type SystemInfo struct {
	Id         peer.ID   `json:"node_id"`
	UUID       string    `json:"uuid"`
	Ip         string    `json:"ip"`
	Network    string    `json:"network"`
	Hostname   string    `json:"node,omitempty"`
	OS         string    `json:"os"`
	Platform   string    `json:"platform"`
	Version    string    `json:"version"`
	Role       string    `json:"role"`
	OnlineUser int       `json:"online_user"`
	Time       time.Time `json:"time"`
}

// NewSystemInfo create method
func NewSystemInfo(ip string, nodeId peer.ID, role string, network string) *SystemInfo {

	var platform, version, oS = platformInformation()
	var host, _ = os.Hostname()

	return &SystemInfo{
		Id:       nodeId,
		Ip:       ip,
		Network:  network,
		Hostname: host,
		Platform: platform,
		Version:  version,
		OS:       oS,
		Role:     role,
	}
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

func (p *SystemInfo) UpdateLoggedInUser() {
	if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
		output, _ := exec.Command("who").Output()
		users := outputToIntUserLinux(string(output))
		p.OnlineUser = users
		return
	}
	if runtime.GOOS == "windows" {
		output, _ := exec.Command("query", "user").Output()
		users := outputToIntUserWin(string(output))
		p.OnlineUser = users
		return
	}
	return
}

func outputToIntUserWin(output string) int {
	var users = 0

	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		words := strings.Fields(line)
		if strings.HasPrefix(words[3], "Active") {
			users++
		}
	}
	err := scanner.Err()
	if err != nil {
		log.Println(err)
		return 0
	}
	return users
}

func outputToIntUserLinux(output string) int {
	var users = 0

	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		users++
	}
	err := scanner.Err()
	if err != nil {
		return 0
	}
	return users
}
