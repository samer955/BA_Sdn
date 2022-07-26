package variables

import "time"

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
	Id        string    `json:"node_id"`
	Ip        string    `json:"ip"`
	Processes []Process `json:"processes"`
	Hostname  string    `json:"host,omitempty"`
	Usage     int       `json:"usage, omitempty"`
	Time      time.Time `json:"time, omitempty"`
}

type Process struct {
	Name        string  `json:"name"`
	Cpu_percent float64 `json:"cpu_percent"`
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
