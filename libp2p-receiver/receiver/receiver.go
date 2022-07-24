package receiver

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/beevik/ntp"
	host2 "github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/p2p/protocol/ping"
	"libp2p-receiver/database"
	"libp2p-receiver/variables"
	"log"
	"time"
)

var db = database.Database()

func ReadTimeMessages(subscribe *pubsub.Subscription, context context.Context, node host2.Host) {
	for {
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
				_, err = db.Exec("INSERT INTO latency(hostname,node_id,ip,latency,time) "+
					"VALUES($1,$2,$3,$4,$5)", peer.Hostname, peer.Id, peer.Ip, latency, now)

				if err != nil {
					log.Fatal(err)
				}
			}
		}
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

// SendPing In order to see if a Peer is alive or not, send a Ping and get an RTT response,
//If the ping return an error, we cannot reach it.
func SendPing(ctx context.Context, node host2.Host, peer peer.AddrInfo) {

	status := variables.PingStatus{
		Source_node: node.ID().Pretty(),
		Target_node: peer.ID.Pretty(),
	}
	//The Ping function return a channel that still open till the context is alive
	ch := ping.Ping(ctx, node, peer.ID)

	for {
		res := <-ch

		if res.Error == nil {
			status.Alive = true
			status.RTT = res.RTT.Milliseconds()
			fmt.Println("pinged", peer.Addrs[0], "in", res.RTT)
		} else {
			status.Alive = false
			status.RTT = 0
			fmt.Println("pinged", peer.Addrs[0], "without success")
			return
		}
		status.Time = TimeFromServer()
		savePingStatus(status)

		//Next Ping in 1 Min
		time.Sleep(10 * time.Second)
	}
}

func savePingStatus(status variables.PingStatus) {

	_, err := db.Exec("INSERT INTO status(source_id,target_id,is_alive,rtt,time)"+
		" VALUES($1,$2,$3,$4,$5)",
		status.Source_node,
		status.Target_node,
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
