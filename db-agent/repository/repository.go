package repository

import (
	"db-agent/variables"
)

type Repository interface {
	Migrate()
	SaveSystemInfo(peer *variables.SystemInfo) *variables.SystemInfo
	SaveRamInfo(ram *variables.Ram) *variables.Ram
	SaveCpuIfo(cpu *variables.Cpu) *variables.Cpu
	SavePingStatus(status *variables.PingStatus) *variables.PingStatus
	SaveTCPstatus(tcpStatus *variables.TCPstatus) *variables.TCPstatus
	SaveBandwidth(data *variables.Bandwidth) *variables.Bandwidth
	GetIpFromNode(node string) string
}
