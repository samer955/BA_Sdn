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
	"github.com/multiformats/go-multiaddr"
	"libp2p-receiver/database"
	"libp2p-receiver/variables"
	"log"
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
				now, err := ntp.Time("time.apple.com")
				if err != nil {
					fmt.Println(err)
				}

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

				log.Printf("Message: <%s> %s", message.Data, message.ReceivedFrom.String())
			}
		}
	}
}

func ReadPingMessage(subscribe *pubsub.Subscription, ctx context.Context, node host2.Host, pingService *ping.PingService) {

	for {
		message, err := subscribe.Next(ctx)
		if err != nil {
			log.Println("cannot read from topic")
		} else {
			if message.ReceivedFrom.String() != node.ID().Pretty() {

				var ping string

				//parse the JSON-encoded data and store the result into cpu
				json.Unmarshal(message.Data, &ping)

				fmt.Println(ping)

				addr, err := multiaddr.NewMultiaddr(ping)
				if err != nil {
					panic(err)
				}

				peer, err := peer.AddrInfoFromP2pAddr(addr)
				if err != nil {
					panic(err)
				}

				fmt.Println(peer)
				fmt.Print(addr)
				//	if err := node.Connect(context.Background(), *peer); err != nil {
				//		panic(err)
				//	}
				//
				fmt.Println("sending 5 ping messages to", addr)
				ch := pingService.Ping(context.Background(), peer.ID)
				for i := 0; i < 5; i++ {
					res := <-ch
					fmt.Println("pinged", addr, "in", res.RTT.Milliseconds())
				}
				//
				//	fmt.Println(ping)
			}
		}

	}
}

func latencyCalculate(actual, source int64) int64 {
	return actual - source
}
