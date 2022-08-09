package service

import (
	"database/sql"
	"db-agent/variables"
	"fmt"
	_ "github.com/lib/pq"
	"io/ioutil"
	"log"
	"time"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "user"
	password = "password"
	dbname   = "database"
)

var myDb *sql.DB

func init() {

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	postgresDb, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		panic(err)
	} else {
		fmt.Println("Successfully connected!")
	}

	tables, err := ioutil.ReadFile("database/migrations/000001_init_schema.up.sql")
	if err != nil {
		log.Fatal(err)
	}

	// Convert []byte to string and print to screen
	text := string(tables)
	postgresDb.Exec(text)

	myDb = postgresDb
}

func SaveSystemMessage(peer *variables.PeerInfo, now time.Time, latency int64) {

	_, err := myDb.Exec("INSERT INTO peer(node_id,uuid,hostname,ip,latency,platform,version,os,online_user,time) "+
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

func SavePingStatus(status variables.PingStatus) {

	_, err := myDb.Exec("INSERT INTO status(uuid,source,target,is_alive,rtt,time)"+
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

func SaveRamInfo(ram *variables.Ram) {

	_, err := myDb.Exec("INSERT INTO ram(uuid,hostname,node_id,ip,usage,time)"+
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

func SaveCpuIfo(cpu *variables.Cpu) {

	_, err := myDb.Exec("INSERT INTO cpu(uuid,hostname,node_id,ip,usage,model,time)"+
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

func SaveTCPstatus(tcpStatus variables.TCPstatus) {
	_, err := myDb.Exec("INSERT INTO tcp(uuid,hostname,ip,queue_size,received,sent,time) "+
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

func SaveThroughput(data *variables.IOData) {
	_, err := myDb.Exec("INSERT INTO throughput(uuid,node_id,ip,total_in,total_out,rate_in,rate_out,hostname,time) "+
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
