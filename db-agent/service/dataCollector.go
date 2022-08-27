package service

import (
	"context"
	"db-agent/repository"
	"db-agent/variables"
	"encoding/json"
	"fmt"
	"github.com/beevik/ntp"
	"github.com/libp2p/go-libp2p-core/host"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"log"
	"time"
)

type dataCollector struct {
	node       host.Host
	repository *repository.PostGresRepo
}

func NewDataCollectorService(node host.Host, repo *repository.PostGresRepo) *dataCollector {
	return &dataCollector{
		node:       node,
		repository: repo,
	}
}

func (collector *dataCollector) ReadSystemInfo(subscribe *pubsub.Subscription, context context.Context) {
	for {
		func() {
			defer handlePanicError()
			message, err := subscribe.Next(context)
			if err != nil {
				log.Println("cannot read from topic")
			} else {
				if message.ReceivedFrom.String() != collector.node.ID().Pretty() {

					peer := new(variables.PeerInfo)

					//Unmarshal the file into the peer struct
					json.Unmarshal(message.Data, peer)

					//Get the actual time from a remote Server
					now := TimeFromServer()

					//Latency is calculated from the time when the peer send the message
					//and when the service reads it (in millis)
					latency := latencyCalculate(now.UnixMilli(), peer.Time.UnixMilli())

					log.Println("latency :", latency)

					//Here we store latency of the peer in the database as well as system information
					collector.repository.SaveSystemMessage(peer, now, latency)

					log.Printf("Message: <%s> %s", message.Data, message.ReceivedFrom.String())
				}
			}
		}()
	}
}

func (collector *dataCollector) ReadRamInformation(subscribe *pubsub.Subscription, ctx context.Context) {
	for {
		func() {
			defer handlePanicError()
			message, err := subscribe.Next(ctx)
			if err != nil {
				log.Println("cannot read from topic")
			} else {
				if message.ReceivedFrom.String() != collector.node.ID().Pretty() {

					ram := new(variables.Ram)

					//parse the JSON-encoded data and store the result into ram
					json.Unmarshal(message.Data, ram)

					//Here we store cpu usage percentage of the peer in the database as well
					//as node_id, ip_address to identify the peer
					collector.repository.SaveRamInfo(ram)

					log.Printf("Message: <%s> %s", message.Data, message.ReceivedFrom.String())
				}
			}
		}()
	}
}

func (collector *dataCollector) ReadCpuInformation(subscribe *pubsub.Subscription, ctx context.Context) {
	for {
		func() {
			defer handlePanicError()
			message, err := subscribe.Next(ctx)
			if err != nil {
				log.Println("cannot read from topic")
			} else {
				if message.ReceivedFrom.String() != collector.node.ID().Pretty() {

					cpu := new(variables.Cpu)

					//parse the JSON-encoded data and store the result into cpu
					json.Unmarshal(message.Data, cpu)

					//Here we store cpu usage percentage of the peer in the database as well
					//as node_id, ip_address to identify the peer
					collector.repository.SaveCpuIfo(cpu)

					log.Printf("Message: <%s> %s", message.Data, message.ReceivedFrom.String())
				}
			}
		}()
	}
}

func (collector *dataCollector) ReadPingStatus(subscribe *pubsub.Subscription, ctx context.Context) {
	for {
		func() {
			defer handlePanicError()
			message, err := subscribe.Next(ctx)
			if err != nil {
				log.Println("cannot read from topic")
			} else {
				if message.ReceivedFrom.String() != collector.node.ID().Pretty() {

					status := new(variables.PingStatus)
					//parse the JSON-encoded data and store the result into cpu
					json.Unmarshal(message.Data, status)

					collector.repository.SavePingStatus(status)
					log.Printf("Message: <%s> %s", message.Data, message.ReceivedFrom.String())
				}
			}
		}()
	}
}

func (collector *dataCollector) ReadTCPstatus(subscribe *pubsub.Subscription, ctx context.Context) {
	for {
		func() {
			defer handlePanicError()
			message, err := subscribe.Next(ctx)
			if err != nil {
				log.Println("cannot read from topic")
			} else {
				if message.ReceivedFrom.String() != collector.node.ID().Pretty() {

					tcpStat := new(variables.TCPstatus)
					//parse the JSON-encoded data and store the result into cpu
					json.Unmarshal(message.Data, tcpStat)

					collector.repository.SaveTCPstatus(tcpStat)
					log.Printf("Message: <%s> %s", message.Data, message.ReceivedFrom.String())
				}
			}
		}()
	}
}

func (collector *dataCollector) ReadBandwidth(subscribe *pubsub.Subscription, ctx context.Context) {
	for {
		func() {
			defer handlePanicError()
			message, err := subscribe.Next(ctx)
			if err != nil {
				log.Println("cannot read from topic")
			} else {
				if message.ReceivedFrom.String() != collector.node.ID().Pretty() {
					log.Printf("Message: <%s> %s", message.Data, message.ReceivedFrom.String())

					bandwidth := new(variables.Bandwidth)
					//parse the JSON-encoded data and store the result into cpu
					json.Unmarshal(message.Data, bandwidth)

					collector.repository.SaveBandwidth(bandwidth)
				}
			}
		}()
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
