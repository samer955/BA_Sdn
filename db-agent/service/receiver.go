package service

import (
	"context"
	"db-agent/node"
	"db-agent/repository"
	"db-agent/variables"
	"encoding/json"
	"fmt"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"log"
	"time"
)

type Receiver struct {
	Node       node.Node
	Repository repository.Repository
}

func (receiver *Receiver) ReadSystemInfo(subscribe *pubsub.Subscription, context context.Context) {
	for {
		func() {
			//defer handlePanicError recovers the state of the program if an error occurs: fundamental if two DB-Agent are used e.g: storing same UUID
			defer handlePanicError()
			message, err := subscribe.Next(context)
			if err != nil {
				log.Println("cannot read from topic")
			} else {
				if message.ReceivedFrom.String() != receiver.Node.Host.ID().Pretty() {
					log.Printf("Message: <%s> %s", message.Data, message.ReceivedFrom.String())

					//Get the actual time
					now := time.Now()
					//empty systemInfo struct, here will be copied the content of the message
					systemInfo := new(variables.SystemInfo)
					//Unmarshal the file into the systemInfo struct
					json.Unmarshal(message.Data, systemInfo)
					//Latency is calculated from the time when the sender-agent send the message
					//and when the service reads it (in ms)
					latency := latencyCalculate(now.UnixMilli(), systemInfo.Time.UnixMilli())

					log.Printf("latency node %s: %d ms\n", message.ReceivedFrom.Pretty(), latency)

					//Storing system info in the db
					receiver.Repository.SaveSystemInfo(systemInfo, now, latency)
				}
			}
		}()
	}
}

func (receiver *Receiver) ReadRamInformation(subscribe *pubsub.Subscription, ctx context.Context) {
	for {
		func() {
			defer handlePanicError()
			message, err := subscribe.Next(ctx)
			if err != nil {
				log.Println("cannot read from topic")
			} else {
				if message.ReceivedFrom.String() != receiver.Node.Host.ID().Pretty() {
					log.Printf("Message: <%s> %s", message.Data, message.ReceivedFrom.String())

					ram := new(variables.Ram)
					//parse the JSON-encoded data and store the result into ram
					json.Unmarshal(message.Data, ram)

					receiver.Repository.SaveRamInfo(ram)
				}
			}
		}()
	}
}

func (receiver *Receiver) ReadCpuInformation(subscribe *pubsub.Subscription, ctx context.Context) {
	for {
		func() {
			defer handlePanicError()
			message, err := subscribe.Next(ctx)
			if err != nil {
				log.Println("cannot read from topic")
			} else {
				if message.ReceivedFrom.String() != receiver.Node.Host.ID().Pretty() {
					log.Printf("Message: <%s> %s", message.Data, message.ReceivedFrom.String())

					cpu := new(variables.Cpu)

					//parse the JSON-encoded data and store the result into cpu
					json.Unmarshal(message.Data, cpu)

					receiver.Repository.SaveCpuIfo(cpu)
				}
			}
		}()
	}
}

func (receiver *Receiver) ReadPingStatus(subscribe *pubsub.Subscription, ctx context.Context) {
	for {
		func() {
			defer handlePanicError()
			message, err := subscribe.Next(ctx)
			if err != nil {
				log.Println("cannot read from topic")
			} else {
				if message.ReceivedFrom.String() != receiver.Node.Host.ID().Pretty() {
					log.Printf("Message: <%s> %s", message.Data, message.ReceivedFrom.String())

					pingStatus := new(variables.PingStatus)
					//parse the JSON-encoded data and store the result into cpu
					json.Unmarshal(message.Data, pingStatus)

					//get the ip source/target from the repository and save it in the pingStatus
					sourceIp := receiver.Repository.GetIpFromNode(pingStatus.Source)
					targetIp := receiver.Repository.GetIpFromNode(pingStatus.Target)
					pingStatus.SourceIp = sourceIp
					pingStatus.TargetIp = targetIp
					receiver.Repository.SavePingStatus(pingStatus)
				}
			}
		}()
	}
}

func (receiver *Receiver) ReadTCPstatus(subscribe *pubsub.Subscription, ctx context.Context) {
	for {
		func() {
			defer handlePanicError()
			message, err := subscribe.Next(ctx)
			if err != nil {
				log.Println("cannot read from topic")
			} else {
				if message.ReceivedFrom.String() != receiver.Node.Host.ID().Pretty() {
					log.Printf("Message: <%s> %s", message.Data, message.ReceivedFrom.String())

					tcpStat := new(variables.TCPstatus)
					//parse the JSON-encoded data and store the result into cpu
					json.Unmarshal(message.Data, tcpStat)

					receiver.Repository.SaveTCPstatus(tcpStat)
				}
			}
		}()
	}
}

func (receiver *Receiver) ReadBandwidth(subscribe *pubsub.Subscription, ctx context.Context) {
	for {
		func() {
			defer handlePanicError()
			message, err := subscribe.Next(ctx)
			if err != nil {
				log.Println("cannot read from topic")
			} else {
				if message.ReceivedFrom.String() != receiver.Node.Host.ID().Pretty() {
					log.Printf("Message: <%s> %s", message.Data, message.ReceivedFrom.String())

					bandwidth := new(variables.Bandwidth)
					//parse the JSON-encoded data and store the result into cpu
					json.Unmarshal(message.Data, bandwidth)

					receiver.Repository.SaveBandwidth(bandwidth)
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
