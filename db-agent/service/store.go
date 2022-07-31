package service

import (
	"db-agent/variables"
	"log"
	"time"
)

func saveSystemMessage(peer *variables.PeerInfo, now time.Time, latency int64) {

	_, err := db.Exec("INSERT INTO peer(node_id,uuid,hostname,ip,latency,platform,version,os,time) "+
		"VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9)",
		peer.Id,
		peer.UUID,
		peer.Hostname,
		peer.Ip,
		latency,
		peer.Platform,
		peer.Version,
		peer.OS,
		now)

	if err != nil {
		log.Println(err)
	}
}

func savePingStatus(status variables.PingStatus) {

	_, err := db.Exec("INSERT INTO status(uuid,source,target,is_alive,rtt,time)"+
		" VALUES($1,$2,$3,$4,$5,$6)",
		status.UUID,
		status.Source,
		status.Target,
		status.Alive,
		status.RTT,
		status.Time,
	)

	if err != nil {
		log.Println(err)
	}
}

func saveRamInfo(ram *variables.Ram) {

	_, err := db.Exec("INSERT INTO ram(uuid,hostname,ip,usage,time)"+
		" VALUES($1,$2,$3,$4,$5,$6)",
		ram.UUID,
		ram.Hostname,
		ram.Id,
		ram.Ip,
		ram.Usage,
		ram.Time)

	if err != nil {
		log.Println(err)
	}
}

func saveCpuIfo(cpu *variables.Cpu) {

	_, err := db.Exec("INSERT INTO cpu(uuid,hostname,node_id,ip,usage,model,time)"+
		" VALUES($1,$2,$3,$4,$5,$6,$7)",

		cpu.Hostname,
		cpu.Id,
		cpu.Ip,
		cpu.Usage,
		cpu.Time)

	if err != nil {
		log.Println(err)
	}
}

func saveTCPstatus(tcpStatus variables.TCPstatus) {
	_, err := db.Exec("INSERT INTO tcp(uuid,hostname,ip,queue_size,received,sent) "+
		"VALUES($1,$2,$3,$4,$5,$6)",
		tcpStatus.UUID,
		tcpStatus.Hostname,
		tcpStatus.Ip,
		tcpStatus.QueueSize,
		tcpStatus.Received,
		tcpStatus.Sent)

	if err != nil {
		log.Println(err)
	}
}
