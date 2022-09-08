package repository

import (
	"database/sql"
	"db-agent/config"
	"db-agent/variables"
	"log"
	"time"
)

type PostGresRepo struct {
	db *sql.DB
}

func NewPostGresRepository(db *sql.DB) *PostGresRepo {
	return &PostGresRepo{db: db}
}

var conf = config.GetConfig()

func (r *PostGresRepo) Migrate() {
	_, err := r.db.Exec(conf.TableSchema)
	if err != nil {
		log.Println("unable to execute migration")
		panic(err)
	}
}

func (r *PostGresRepo) SaveSystemInfo(peer *variables.SystemInfo, now time.Time, latency int64) {
	_, err := r.db.Exec("INSERT INTO peer(node_id,uuid,hostname,ip,latency,platform,version,os,online_user,time,role,network) "+
		"VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)",
		peer.Id.Pretty(),
		peer.UUID,
		peer.Hostname,
		peer.Ip,
		latency,
		peer.Platform,
		peer.Version,
		peer.OS,
		peer.OnlineUser,
		now,
		peer.Role,
		peer.Network)

	if err != nil {
		panic(err)
	}
}

func (r *PostGresRepo) SavePingStatus(status *variables.PingStatus) {

	_, err := r.db.Exec("INSERT INTO status(uuid,source,target,is_alive,rtt,time,source_ip,target_ip)"+
		" VALUES($1,$2,$3,$4,$5,$6,$7,$8)",
		status.UUID,
		status.Source,
		status.Target,
		status.Alive,
		status.RTT,
		status.Time,
		status.SourceIp,
		status.TargetIp)

	if err != nil {
		panic(err)
	}
}

func (r *PostGresRepo) SaveRamInfo(ram *variables.Ram) {

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
}

func (r *PostGresRepo) SaveCpuIfo(cpu *variables.Cpu) {

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
}

func (r *PostGresRepo) SaveTCPstatus(tcpStatus *variables.TCPstatus) {
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
}

func (r *PostGresRepo) SaveBandwidth(band *variables.Bandwidth) {
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
}

func (r *PostGresRepo) GetIpFromNode(node string) string {
	ip := ""
	sqlStatement := `SELECT ip FROM peer WHERE node_id=$1;`
	// Query for a value based on a single row.
	r.db.QueryRow(sqlStatement, node).Scan(&ip)
	return ip
}
