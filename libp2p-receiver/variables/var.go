package variables

import "time"

type PeerInfo struct {
	Id       string    `json:"node_id"`
	Ip       string    `json:"ip"`
	Hostname string    `json:"host,omitempty"`
	Time     time.Time `json:"time"`
}

type Cpu struct {
	Id       string    `json:"node_id"`
	Ip       string    `json:"ip"`
	Hostname string    `json:"host,omitempty"`
	Usage    int       `json:"usage, omitempty"`
	Time     time.Time `json:"time, omitempty"`
}

type Ram struct {
	Id       string    `json:"node_id"`
	Ip       string    `json:"ip"`
	Hostname string    `json:"host,omitempty"`
	Usage    int       `json:"usage, omitempty"`
	Time     time.Time `json:"time, omitempty"`
}
