package service

import (
	"context"
	"db-agent/database"
	"db-agent/variables"
	"encoding/json"
	"fmt"
	"github.com/beevik/ntp"
	"github.com/google/uuid"
	host2 "github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/metrics"
	"github.com/libp2p/go-libp2p-core/peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"log"
	"time"
)

var db = database.Database()
var nodeIpMap = make(map[string]string)

func ReadSystemInfo(subscribe *pubsub.Subscription, context context.Context, node host2.Host) {
	for {
		func() {
			defer handlePanicError()
			message, err := subscribe.Next(context)
			if err != nil {
				log.Println("cannot read from topic")
			} else {
				if message.ReceivedFrom.String() != node.ID().Pretty() {

					peer := new(variables.PeerInfo)

					//Unmarshal the file into the peer struct
					json.Unmarshal(message.Data, peer)

					nodeIpMap[peer.Id] = peer.Ip

					fmt.Println(nodeIpMap)

					//Get the actual time from a remote Server
					now := TimeFromServer()

					//Latency is calculated from the time when the peer send the message
					//and when the service reads it (in millis)
					latency := latencyCalculate(now.UnixMilli(), peer.Time.UnixMilli())

					log.Println("latency :", latency)

					//Here we store latency of the peer in the database as well as system information
					saveSystemMessage(peer, now, latency)

					log.Printf("Message: <%s> %s", message.Data, message.ReceivedFrom.String())
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
				saveRamInfo(ram)

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
				saveCpuIfo(cpu)

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

func ReadTCPstatus(subscribe *pubsub.Subscription, ctx context.Context, node host2.Host) {
	for {
		message, err := subscribe.Next(ctx)
		if err != nil {
			log.Println("cannot read from topic")
		} else {
			if message.ReceivedFrom.String() != node.ID().Pretty() {

				tcpStat := new(variables.TCPstatus)
				//parse the JSON-encoded data and store the result into cpu
				json.Unmarshal(message.Data, tcpStat)

				saveTCPstatus(*tcpStat)
				log.Printf("Message: <%s> %s", message.Data, message.ReceivedFrom.String())
			}
		}
	}
}

func ReadBandwidth(counter *metrics.BandwidthCounter, peerlist *[]peer.AddrInfo) {

	ioData := new(variables.IOData)

	for {
		if len(*peerlist) == 0 {
			continue
		}
		mapPeer := counter.GetBandwidthByPeer()
		now := TimeFromServer()

		for key, value := range mapPeer {
			ioData.NodeID = key.Pretty()
			ioData.TotalIn = value.TotalIn
			ioData.TotalOut = value.TotalOut
			ioData.RateIn = int(value.RateIn)
			ioData.RateOut = int(value.TotalOut)
			ioData.UUID = uuid.New().String()
			ioData.Time = now

			fmt.Println(ioData)
		}
		time.Sleep(60 * time.Second)
	}
}

//TimeFromServer get the actual time from a remote server using the ntp Protocol
//The purpose is to synchronize the time between the VMs to avoid problems
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
