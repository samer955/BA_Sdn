package service

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/p2p/protocol/ping"
	"libp2p-sender/components"
	"libp2p-sender/subscriber"
	"time"
)

//SendPeerInfo will send system Information of this Peer and periodically a timestamp
//in order to calculate the latency in ms between service and service
func SendPeerInfo(topic *pubsub.Topic, context context.Context, peer *components.PeerInfo, list *[]peer.AddrInfo) {
	for {
		if len(*list) == 0 {
			continue
		}
		sendPeerInfo(topic, context, peer)

		//wait 20 seconds before send another timestamp
		time.Sleep(20 * time.Second)
	}
}

func sendPeerInfo(topic *pubsub.Topic, context context.Context, peer *components.PeerInfo) {
	fmt.Println("sending time")

	//Set the time when the message is sent
	peer.UUID = uuid.New().String()
	peer.OnlineUser = components.GetNumberOfOnlineUser()
	peer.Time = components.TimeFromServer()

	err := subscriber.Publish(peer, context, topic)
	if err != nil {
		fmt.Println("Error publishing content ", err.Error())
	}
}

// SendCpuInfo function will send periodically information about the CPU
func SendCpuInfo(topic *pubsub.Topic, context context.Context, cpu *components.Cpu, peers *[]peer.AddrInfo) {
	for {
		if len(*peers) == 0 {
			continue
		}
		sendCpuInfo(topic, context, cpu)

		time.Sleep(15 * time.Second)
	}
}

func sendCpuInfo(topic *pubsub.Topic, context context.Context, cpu *components.Cpu) {
	fmt.Println("sending cpu")

	cpu.UUID = uuid.New().String()
	//Update every 10s CPU Usages in %
	cpu.UpdateCpuPercentage()

	//publish the cpu data
	err := subscriber.Publish(cpu, context, topic)

	if err != nil {
		fmt.Println("Error publishing content ", err.Error())
	}
}

// SendRamInfo function will send periodically information about the actual RAM Percentage
func SendRamInfo(topic *pubsub.Topic, context context.Context, ram *components.Ram, peers *[]peer.AddrInfo) {
	for {
		if len(*peers) == 0 {
			continue
		}
		sendRamInfo(topic, context, ram)
		time.Sleep(5 * time.Second)
	}
}

func sendRamInfo(topic *pubsub.Topic, context context.Context, ram *components.Ram) {
	fmt.Println("sending ram")

	ram.UUID = uuid.New().String()
	//Update every 10s RAM usages in %
	ram.UpdateRamPercentage()

	//publish the ram data
	err := subscriber.Publish(ram, context, topic)
	if err != nil {
		fmt.Println("Error publishing content ", err.Error())
	}
}

//SendPing function send a Ping every 60s to the discovered Peer
//This function is used in order to reach the others nodes in an active way getting
//a bool if the ping was successfully (true if a node is reachable, false if not) and an RTT in ms
func SendPing(ctx context.Context, host host.Host, target peer.AddrInfo, topic *pubsub.Topic) {

	var pingDeadline = 10

	status := components.NewPingStatus(host.ID().Pretty(), target.ID.Pretty())
	//The Ping function return a channel that still open till the context is alive
	ch := ping.Ping(ctx, host, target.ID)

	for {
		//after 10 negative Ping stop to ping the Peer
		if pingDeadline == 0 {
			fmt.Printf("Stopped ping from %s to %s\n", status.Source, status.Target)
			return
		}
		res := <-ch

		status.SetPingStatus(res, &pingDeadline)

		//publish the status of the Ping in the topic
		subscriber.Publish(status, ctx, topic)

		//Next Ping in 1 Min
		time.Sleep(5 * time.Second)
	}
}

func SendTCPstatus(topic *pubsub.Topic, context context.Context, tcpIfo *components.TCPstatus, peers *[]peer.AddrInfo) {
	for {
		if len(*peers) == 0 {
			continue
		}
		sendTCPstatus(topic, context, tcpIfo)
		time.Sleep(15 * time.Second)
	}
}

func sendTCPstatus(topic *pubsub.Topic, context context.Context, tcpIfo *components.TCPstatus) {
	fmt.Println("sending tcp Queue size")

	tcpIfo.UUID = uuid.New().String()
	tcpIfo.QueueSize = components.TcpQueueSize()
	received, sent := components.TcpSegmentsNumber()

	tcpIfo.Received = received
	tcpIfo.Sent = sent
	tcpIfo.Time = components.TimeFromServer()

	err := subscriber.Publish(tcpIfo, context, topic)
	if err != nil {
		fmt.Println("Error publishing content ", err.Error())
	}
}
