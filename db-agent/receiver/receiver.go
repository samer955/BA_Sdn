package receiver

import (
	"context"
	"db-agent/database"
	"db-agent/variables"
	"encoding/json"
	"fmt"
	"github.com/beevik/ntp"
	host2 "github.com/libp2p/go-libp2p-core/host"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"log"
	"time"
)

var db = database.Database()

func ReadTimeMessages(subscribe *pubsub.Subscription, context context.Context, node host2.Host) {
	for {
		func() {
			defer handlePanicError()
			message, err := subscribe.Next(context)
			if err != nil {
				log.Println("cannot read from topic")
			} else {
				if message.ReceivedFrom.String() != node.ID().Pretty() {
					log.Printf("Message: <%s> %s", message.Data, message.ReceivedFrom.String())

					peer := new(variables.PeerInfo)

					//Unmarshal the file into the peer struct
					json.Unmarshal(message.Data, peer)

					//Get the actual time from a Server
					now := TimeFromServer()

					//Latency is calculated from the time when the peer send the message
					//and when the receiver reads it (in millis)
					latency := latencyCalculate(now.UnixMilli(), peer.Time.UnixMilli())

					log.Println("latency :", latency)

					//Here we store latency of the peer in the database as well as node_id, ip_address
					//in order to identify it
					saveTimeMessage(peer, now, latency)
				}
			}
		}()
	}
}

func ReadRamInformation(subscribe *pubsub.Subscription, ctx context.Context, node host2.Host) {
	for {
		message, err := subscribe.Next(ctx)
		if err != nil {
			log.Println("cannot read from topic")
		} else {
			if message.ReceivedFrom.String() != node.ID().Pretty() {

				ram := new(variables.Ram)

				//parse the JSON-encoded data and store the result into ram
				json.Unmarshal(message.Data, ram)

				//Here we store cpu usage percentage of the peer in the database as well
				//as node_id, ip_address to identify the peer
				_, err = db.Exec("INSERT INTO ram(hostname,node_id,ip,usage,time)"+
					" VALUES($1,$2,$3,$4,$5)",
					ram.Hostname,
					ram.Id,
					ram.Ip,
					ram.Usage,
					ram.Time)

				if err != nil {
					log.Fatal(err)
				}

				log.Printf("Message: <%s> %s", message.Data, message.ReceivedFrom.String())
			}
		}
	}
}

func ReadCpuInformation(subscribe *pubsub.Subscription, ctx context.Context, node host2.Host) {
	for {
		message, err := subscribe.Next(ctx)
		if err != nil {
			log.Println("cannot read from topic")
		} else {
			if message.ReceivedFrom.String() != node.ID().Pretty() {

				cpu := new(variables.Cpu)

				//parse the JSON-encoded data and store the result into cpu
				json.Unmarshal(message.Data, cpu)

				//Here we store cpu usage percentage of the peer in the database as well
				//as node_id, ip_address to identify the peer
				_, err = db.Exec("INSERT INTO cpu(hostname,node_id,ip,usage,time)"+
					" VALUES($1,$2,$3,$4,$5)",
					cpu.Hostname,
					cpu.Id,
					cpu.Ip,
					cpu.Usage,
					cpu.Time)

				if err != nil {
					log.Fatal(err)
				}

				if len(cpu.Processes) != 0 {
					for _, proc := range cpu.Processes {

						_, err = db.Exec("INSERT INTO process(name,cpu,hostname,ip,time)"+
							" VALUES($1,$2,$3,$4,$5)",
							proc.Name,
							proc.Cpu_percent,
							cpu.Hostname,
							cpu.Ip,
							cpu.Time)

						if err != nil {
							log.Fatal(err)
						}

					}
				}

				log.Printf("Message: <%s> %s", message.Data, message.ReceivedFrom.String())
			}
		}
	}
}

func ReadPingStatus(subscribe *pubsub.Subscription, ctx context.Context, node host2.Host) {

	for {
		message, err := subscribe.Next(ctx)
		if err != nil {
			log.Println("cannot read from topic")
		} else {
			if message.ReceivedFrom.String() != node.ID().Pretty() {

				status := new(variables.PingStatus)
				//parse the JSON-encoded data and store the result into cpu
				json.Unmarshal(message.Data, status)

				savePingStatus(*status)
				log.Printf("Message: <%s> %s", message.Data, message.ReceivedFrom.String())
			}
		}

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
		log.Fatal(err)
	}
}

func TimeFromServer() time.Time {
	now, err := ntp.Time("time.apple.com")
	if err != nil {
		fmt.Println(err)
	}
	return now
}

func latencyCalculate(actual, source int64) int64 {
	return actual - source
}

func handlePanicError() {
	if r := recover(); r != nil {
		fmt.Println("Recovered. Error:\n", r)
	}
}

func saveTimeMessage(peer *variables.PeerInfo, now time.Time, latency int64) {

	_, err := db.Exec("INSERT INTO peer(uuid,hostname,ip,latency,platform,version,os,time) "+
		"VALUES($1,$2,$3,$4,$5,$6,$7,$8)",
		peer.UUID,
		peer.Hostname,
		peer.Ip,
		latency,
		peer.Platform,
		peer.Version,
		peer.OS,
		now)

	if err != nil {
		panic(err)
	}

}
