package metrics

import (
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/shirou/gopsutil/v3/host"
	"net"
	"os"
	"time"
)

type PeerInfo struct {
	Id         peer.ID   `json:"node_id"`
	UUID       string    `json:"uuid"`
	Ip         string    `json:"ip"`
	Network    string    `json:"network"`
	Hostname   string    `json:"host,omitempty"`
	OS         string    `json:"os"`
	Platform   string    `json:"platform"`
	Version    string    `json:"version"`
	Role       string    `json:"role"`
	OnlineUser int       `json:"online_user"`
	Time       time.Time `json:"time"`
}

// NewPeerInfo create method
func NewPeerInfo(ip string, nodeId peer.ID, role string, network string) *PeerInfo {

	var platform, version, oS = platformInformation()
	var host, _ = os.Hostname()

	return &PeerInfo{
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

// LocalIP get the host machine local IP address
func LocalIP() string {
	ifaces, err := net.Interfaces()
	if err != nil {
		return ""
	}
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			return ""
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip.IsPrivate() {
				return ip.String()
			}
		}
	}
	return ""
}
