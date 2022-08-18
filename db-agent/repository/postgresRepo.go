package repository

import (
	"database/sql"
	"db-agent/variables"
	"io/ioutil"
	"log"
	"time"
)

type PostGresRepo struct {
	db *sql.DB
}

func NewPostGresRepository(db *sql.DB) *PostGresRepo {
	return &PostGresRepo{db: db}
}

//Create tables from the migrations file
func (r *PostGresRepo) Migrate() {
	tables, err := ioutil.ReadFile("repository/migrations/000001_init_schema.up.sql")
	if err != nil {
		log.Fatal(err)
	}
	text := string(tables)
	r.db.Exec(text)
}

func (r *PostGresRepo) SaveSystemMessage(peer *variables.PeerInfo, now time.Time, latency int64) {
	_, err := r.db.Exec("INSERT INTO peer(node_id,uuid,hostname,ip,latency,platform,version,os,online_user,time) "+
		"VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)",
		peer.Id.Pretty(),
		peer.UUID,
		peer.Hostname,
		peer.Ip,
		latency,
		peer.Platform,
		peer.Version,
		peer.OS,
		peer.OnlineUser,
		now)

	if err != nil {
		panic(err)
	}
}

func (r *PostGresRepo) SavePingStatus(status *variables.PingStatus) {

	_, err := r.db.Exec("INSERT INTO status(uuid,source,target,is_alive,rtt,time)"+
		" VALUES($1,$2,$3,$4,$5,$6)",
		status.UUID,
		status.Source,
		status.Target,
		status.Alive,
		status.RTT,
		status.Time,
	)

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

func (r *PostGresRepo) SaveThroughput(data *variables.IOData) {
	_, err := r.db.Exec("INSERT INTO throughput(uuid,node_id,ip,total_in,total_out,rate_in,rate_out,hostname,time) "+
		"VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9)",
		data.UUID,
		data.NodeID,
		data.Ip,
		data.TotalIn,
		data.TotalOut,
		data.RateIn,
		data.RateOut,
		data.Hostname,
		data.Time)

	if err != nil {
		panic(err)
	}
}
