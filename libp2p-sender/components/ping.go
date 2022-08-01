package components

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p/p2p/protocol/ping"
	"time"
)

type PingStatus struct {
	UUID   string    `json:"uuid"`
	Source string    `json:"source"`
	Target string    `json:"target"`
	Alive  bool      `json:"alive"`
	RTT    int64     `json:"rtt"`
	Time   time.Time `json:"time"`
}

func CheckPingStatus(res ping.Result, status PingStatus, peer peer.AddrInfo) {
	if res.Error == nil {
		status.Alive = true
		status.RTT = res.RTT.Milliseconds()
		fmt.Println("pinged", peer.Addrs[0], "in", res.RTT)
	} else {
		status.Alive = false
		status.RTT = 0
		fmt.Println("pinged", peer.Addrs[0], "without success", res.Error)
	}
	status.Time = TimeFromServer()
	status.UUID = uuid.New().String()
}
