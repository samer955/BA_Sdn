package components

import (
	"fmt"
	"github.com/beevik/ntp"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/shirou/gopsutil/v3/host"
	"os"
	"time"
)

type PeerInfo struct {
	Id         peer.ID   `json:"node_id"`
	UUID       string    `json:"uuid"`
	Ip         string    `json:"ip"`
	Hostname   string    `json:"host,omitempty"`
	OS         string    `json:"os"`
	Platform   string    `json:"platform"`
	Version    string    `json:"version"`
	Role       string    `json:"role"`
	OnlineUser int       `json:"online_user"`
	Time       time.Time `json:"time"`
}

// NewPeerInfo create method
func NewPeerInfo(ip string, nodeId peer.ID, role string) *PeerInfo {

	var platform, version, oS = platformInformation()
	var host, _ = os.Hostname()

	return &PeerInfo{
		Id:       nodeId,
		Ip:       ip,
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

//TimeFromServer get the actual time from a remote server using the ntp Protocol
//The purpose is to synchronize the time between the VMs to avoid problems
func TimeFromServer() time.Time {
	now, err := ntp.Time("time.apple.com")
	if err != nil {
		fmt.Println(err)
	}
	return now
}
