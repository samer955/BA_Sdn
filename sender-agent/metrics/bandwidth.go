package metrics

import "time"

//incoming and outgoing data transferred by the local peer.
type Bandwidth struct {
	UUID     string    `json:"uuid"`
	Id       string    `json:"id"`
	Source   string    `json:"source"`
	Target   string    `json:"target"`
	TotalIn  int64     `json:"total_in"`
	TotalOut int64     `json:"total_out"`
	RateIn   int       `json:"rate_in"`
	RateOut  int       `json:"rate_out"`
	Time     time.Time `json:"time"`
}

func NewBandWidth(ip string, nodeId string) *Bandwidth {

	return &Bandwidth{
		Id:     nodeId,
		Source: ip,
	}
}
