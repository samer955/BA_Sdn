package main

import (
	"context"
	"fmt"
	_ "github.com/lib/pq"
	"os"
	"os/signal"
	"sender-agent/config"
	"sender-agent/discovery"
	"sender-agent/node"
	"sender-agent/service"
	"sender-agent/subscriber"
	"syscall"
)

func main() {

	//containing the name of the mdns service and topic names
	const (
		discoveryName = "discoveryRoom"
		roomPing      = "ping"
		roomSystem    = "system"
		roomCpu       = "cpu"
		roomRam       = "ram"
		roomTcp       = "tcp"
		roomBand      = "bandwidth"
	)

	//create a new Node based on Libp2p
	var Node node.Node
	Node.StartNode()

	context := context.Background()

	//return a new pubsub Service using the GossipSub router
	pubsub := subscriber.NewPubSubService(context, Node)

	//Join and Subscribe on different topics
	systemTopic := pubsub.JoinTopic(roomSystem)
	systemSubscr := pubsub.Subscribe(systemTopic)

	pingTopic := pubsub.JoinTopic(roomPing)
	_ = pubsub.Subscribe(pingTopic)

	discovery.SetPingTopic(pingTopic)

	cpuTopic := pubsub.JoinTopic(roomCpu)
	_ = pubsub.Subscribe(cpuTopic)

	ramTopic := pubsub.JoinTopic(roomRam)
	_ = pubsub.Subscribe(ramTopic)

	tcpTopic := pubsub.JoinTopic(roomTcp)
	_ = pubsub.Subscribe(tcpTopic)

	bandTopic := pubsub.JoinTopic(roomBand)
	_ = pubsub.Subscribe(bandTopic)

	// setup local mDNS discovery. If an error occurs try again in 60 seconds
	discovery.SetupDiscovery(Node, discoveryName)

	//creating a new sender to send the metrics
	sender := service.Sender{Node: Node, Frequency: config.GetConfig().Frequency}

	//send Peer-System-Information in a separated thread
	go sender.SendPeerInfo(systemTopic, context)
	//send CPU information in a separated thread
	go sender.SendCpuInfo(cpuTopic, context)
	//send RAM information in a separated thread
	go sender.SendRamInfo(ramTopic, context)
	//send tcp status in a separated thread
	go sender.SendTCPstatus(tcpTopic, context)
	//Get Bandwidth between local Peer and other connected Peers in a separated thread
	go sender.GetBandWidthForActivePeer(systemSubscr, context, bandTopic)

	//Run the program till its stopped
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	fmt.Println("Received signal, shutting down...")
}
