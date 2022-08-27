package repository

import (
	"db-agent/variables"
	"time"
)

type Repository interface {
	Migrate()
	SaveSystemMessage(peer *variables.PeerInfo, now time.Time, latency int64)
	SaveRamInfo(ram *variables.Ram)
	SaveCpuIfo(cpu *variables.Cpu)
	SavePingStatus(status *variables.PingStatus)
	SaveTCPstatus(tcpStatus *variables.TCPstatus)
	SaveBandwidth(data *variables.Bandwidth)
}
