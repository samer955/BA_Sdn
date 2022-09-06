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

type receiver struct {
	node       host.Host
	repository *repository.PostGresRepo
}

func NewDataCollectorService(node host.Host, repo *repository.PostGresRepo) *receiver {
	return &receiver{
		node:       node,
		repository: repo,
	}
}

func (receiver *receiver) ReadSystemInfo(subscribe *pubsub.Subscription, context context.Context) {
	for {
		func() {
			defer handlePanicError()
			message, err := subscribe.Next(context)
			if err != nil {
				log.Println("cannot read from topic")
			} else {
				if message.ReceivedFrom.String() != receiver.node.ID().Pretty() {
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
					receiver.repository.SaveSystemMessage(peer, now, latency)
				}
			}
		}()
	}
}

func (receiver *receiver) ReadRamInformation(subscribe *pubsub.Subscription, ctx context.Context) {
	for {
		func() {
			defer handlePanicError()
			message, err := subscribe.Next(ctx)
			if err != nil {
				log.Println("cannot read from topic")
			} else {
				if message.ReceivedFrom.String() != receiver.node.ID().Pretty() {
					log.Printf("Message: <%s> %s", message.Data, message.ReceivedFrom.String())

					ram := new(variables.Ram)

					//parse the JSON-encoded data and store the result into ram
					json.Unmarshal(message.Data, ram)

					//Here we store cpu usage percentage of the peer in the database as well
					//as node_id, ip_address to identify the peer
					receiver.repository.SaveRamInfo(ram)
				}
			}
		}()
	}
}

func (receiver *receiver) ReadCpuInformation(subscribe *pubsub.Subscription, ctx context.Context) {
	for {
		func() {
			defer handlePanicError()
			message, err := subscribe.Next(ctx)
			if err != nil {
				log.Println("cannot read from topic")
			} else {
				if message.ReceivedFrom.String() != receiver.node.ID().Pretty() {
					log.Printf("Message: <%s> %s", message.Data, message.ReceivedFrom.String())

					cpu := new(variables.Cpu)

					//parse the JSON-encoded data and store the result into cpu
					json.Unmarshal(message.Data, cpu)

					//Here we store cpu usage percentage of the peer in the database as well
					//as node_id, ip_address to identify the peer
					receiver.repository.SaveCpuIfo(cpu)
				}
			}
		}()
	}
}

func (receiver *receiver) ReadPingStatus(subscribe *pubsub.Subscription, ctx context.Context) {
	for {
		func() {
			defer handlePanicError()
			message, err := subscribe.Next(ctx)
			if err != nil {
				log.Println("cannot read from topic")
			} else {
				if message.ReceivedFrom.String() != receiver.node.ID().Pretty() {
					log.Printf("Message: <%s> %s", message.Data, message.ReceivedFrom.String())

					status := new(variables.PingStatus)
					//parse the JSON-encoded data and store the result into cpu
					json.Unmarshal(message.Data, status)

					sourceIp := receiver.repository.GetIpFromNode(status.Source)
					targetIp := receiver.repository.GetIpFromNode(status.Target)
					status.SourceIp = sourceIp
					status.TargetIp = targetIp
					receiver.repository.SavePingStatus(status)
				}
			}
		}()
	}
}

func (receiver *receiver) ReadTCPstatus(subscribe *pubsub.Subscription, ctx context.Context) {
	for {
		func() {
			defer handlePanicError()
			message, err := subscribe.Next(ctx)
			if err != nil {
				log.Println("cannot read from topic")
			} else {
				if message.ReceivedFrom.String() != receiver.node.ID().Pretty() {
					log.Printf("Message: <%s> %s", message.Data, message.ReceivedFrom.String())

					tcpStat := new(variables.TCPstatus)
					//parse the JSON-encoded data and store the result into cpu
					json.Unmarshal(message.Data, tcpStat)

					receiver.repository.SaveTCPstatus(tcpStat)
				}
			}
		}()
	}
}

func (receiver *receiver) ReadBandwidth(subscribe *pubsub.Subscription, ctx context.Context) {
	for {
		func() {
			defer handlePanicError()
			message, err := subscribe.Next(ctx)
			if err != nil {
				log.Println("cannot read from topic")
			} else {
				if message.ReceivedFrom.String() != receiver.node.ID().Pretty() {
					log.Printf("Message: <%s> %s", message.Data, message.ReceivedFrom.String())

					bandwidth := new(variables.Bandwidth)
					//parse the JSON-encoded data and store the result into cpu
					json.Unmarshal(message.Data, bandwidth)

					receiver.repository.SaveBandwidth(bandwidth)
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
