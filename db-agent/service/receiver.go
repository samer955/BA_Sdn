package service

import (
	"context"
	"db-agent/repository"
	"db-agent/variables"
	"encoding/json"
	"fmt"
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
					log.Printf("Message: <%s> %s", message.Data, message.ReceivedFrom.String())

					peer := new(variables.PeerInfo)
					//Unmarshal the file into the peer struct
					json.Unmarshal(message.Data, peer)

					//Get the actual time
					now := time.Now()

					//Latency is calculated from the time when the peer send the message
					//and when the service reads it (in millis)
					latency := latencyCalculate(now.UnixMilli(), peer.Time.UnixMilli())

					log.Println("latency :", latency)

					//Here we store latency of the peer in the database as well as system information
					collector.repository.SaveSystemMessage(peer, now, latency)
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
					log.Printf("Message: <%s> %s", message.Data, message.ReceivedFrom.String())

					ram := new(variables.Ram)

					//parse the JSON-encoded data and store the result into ram
					json.Unmarshal(message.Data, ram)

					//Here we store cpu usage percentage of the peer in the database as well
					//as node_id, ip_address to identify the peer
					collector.repository.SaveRamInfo(ram)
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
					log.Printf("Message: <%s> %s", message.Data, message.ReceivedFrom.String())

					cpu := new(variables.Cpu)

					//parse the JSON-encoded data and store the result into cpu
					json.Unmarshal(message.Data, cpu)

					//Here we store cpu usage percentage of the peer in the database as well
					//as node_id, ip_address to identify the peer
					collector.repository.SaveCpuIfo(cpu)
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
					log.Printf("Message: <%s> %s", message.Data, message.ReceivedFrom.String())

					status := new(variables.PingStatus)
					//parse the JSON-encoded data and store the result into cpu
					json.Unmarshal(message.Data, status)

					sourceIp := collector.repository.GetIpFromNode(status.Source)
					targetIp := collector.repository.GetIpFromNode(status.Target)
					status.SourceIp = sourceIp
					status.TargetIp = targetIp
					collector.repository.SavePingStatus(status)
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
					log.Printf("Message: <%s> %s", message.Data, message.ReceivedFrom.String())

					tcpStat := new(variables.TCPstatus)
					//parse the JSON-encoded data and store the result into cpu
					json.Unmarshal(message.Data, tcpStat)

					collector.repository.SaveTCPstatus(tcpStat)
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

func latencyCalculate(actual, source int64) int64 {
	return actual - source
}

func handlePanicError() {
	if r := recover(); r != nil {
		fmt.Println("Recovered. Error:\n", r)
	}
}
