package variables

import "time"

type PeerInfo struct {
	Id       string    `json:"node_id"`
	UUID     string    `json:"uuid"`
	Ip       string    `json:"ip"`
	Hostname string    `json:"host"`
	OS       string    `json:"os"`
	Platform string    `json:"platform"`
	Version  string    `json:"version"`
	Time     time.Time `json:"time"`
}

type Cpu struct {
	Id       string    `json:"node_id"`
	UUID     string    `json:"uuid"`
	Ip       string    `json:"ip"`
	Hostname string    `json:"host"`
	Model    string    `json:"model"`
	Usage    int       `json:"usage"`
	Time     time.Time `json:"time"`
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

//incoming and outgoing data transferred by the local peer.
type IOData struct {
	NodeID   string
	UUID     string
	TotalIn  int64
	TotalOut int64
	RateIn   int
	RateOut  int
	Time     time.Time
	IP       string
}

//Queue Size = number of open TCP connections
//Received = number of segments received
//Sent = number of segments sent
type TCPstatus struct {
	UUID      string `json:"uuid"`
	Hostname  string `json:"hostname"`
	Ip        string `json:"ip"`
	QueueSize int    `json:"queue_size"`
	Received  int    `json:"received"`
	Sent      int    `json:"sent"`
}
