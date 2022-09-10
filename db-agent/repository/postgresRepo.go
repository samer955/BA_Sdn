package repository

import (
	"database/sql"
	"db-agent/config"
	"db-agent/variables"
	"log"
)

type PostGresRepo struct {
	db *sql.DB
}

func NewPostGresRepository(db *sql.DB) *PostGresRepo {
	return &PostGresRepo{db: db}
}

func (r *PostGresRepo) Migrate() {

	config := config.GetConfig()

	_, err := r.db.Exec(config.TableSchema)
	if err != nil {
		log.Println("unable to execute migration")
		panic(err)
	}

}

func (r *PostGresRepo) SaveSystemInfo(system *variables.SystemInfo) *variables.SystemInfo {

	_, err := r.db.Exec("INSERT INTO system(node_id,uuid,hostname,ip,latency,platform,version,os,online_user,time,role,network) "+
		"VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)",
		system.Id.Pretty(),
		system.UUID,
		system.Hostname,
		system.Ip,
		system.Latency,
		system.Platform,
		system.Version,
		system.OS,
		system.OnlineUser,
		system.Time,
		system.Role,
		system.Network)

	if err != nil {
		panic(err)
	}
	return system

}

func (r *PostGresRepo) SavePingStatus(ping *variables.PingStatus) *variables.PingStatus {

	_, err := r.db.Exec("INSERT INTO ping(uuid,source,target,is_alive,rtt,time,source_ip,target_ip)"+
		" VALUES($1,$2,$3,$4,$5,$6,$7,$8)",
		ping.UUID,
		ping.Source,
		ping.Target,
		ping.Alive,
		ping.RTT,
		ping.Time,
		ping.SourceIp,
		ping.TargetIp)

	if err != nil {
		panic(err)
	}
	return ping

}

func (r *PostGresRepo) SaveRamInfo(ram *variables.Ram) *variables.Ram {

	_, err := r.db.Exec("INSERT INTO ram(uuid,hostname,node_id,ip,usage,time)"+
		" VALUES($1,$2,$3,$4,$5,$6)",
		ram.UUID,
		ram.Hostname,
		ram.Id,
		ram.Ip,
		ram.Usage,
		ram.Time)

	if err != nil {
		panic(err)
	}
	return ram

}

func (r *PostGresRepo) SaveCpuIfo(cpu *variables.Cpu) *variables.Cpu {

	_, err := r.db.Exec("INSERT INTO cpu(uuid,hostname,node_id,ip,usage,model,time)"+
		" VALUES($1,$2,$3,$4,$5,$6,$7)",
		cpu.UUID,
		cpu.Hostname,
		cpu.Id,
		cpu.Ip,
		cpu.Usage,
		cpu.Model,
		cpu.Time)

	if err != nil {
		panic(err)
	}
	return cpu
}

func (r *PostGresRepo) SaveTCPstatus(tcpStatus *variables.TCPstatus) *variables.TCPstatus {

	_, err := r.db.Exec("INSERT INTO tcp(uuid,hostname,ip,queue_size,received,sent,time) "+
		"VALUES($1,$2,$3,$4,$5,$6,$7)",
		tcpStatus.UUID,
		tcpStatus.Hostname,
		tcpStatus.Ip,
		tcpStatus.QueueSize,
		tcpStatus.Received,
		tcpStatus.Sent,
		tcpStatus.Time)

	if err != nil {
		panic(err)
	}
	return tcpStatus

}

func (r *PostGresRepo) SaveBandwidth(band *variables.Bandwidth) *variables.Bandwidth {

	_, err := r.db.Exec("INSERT INTO bandwidth(uuid,id,source,target,total_in,total_out,rate_in,rate_out,time) "+
		"VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9)",
		band.UUID,
		band.Id,
		band.Source,
		band.Target,
		band.TotalIn,
		band.TotalOut,
		band.RateIn,
		band.RateOut,
		band.Time)

	if err != nil {
		panic(err)
	}
	return band

}

func (r *PostGresRepo) GetIpFromNode(node string) string {

	ip := ""
	sqlStatement := `SELECT ip FROM system WHERE node_id=$1;`
	r.db.QueryRow(sqlStatement, node).Scan(&ip)
	return ip

}
