package main

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/metrics"
	"log"
	"os"
	"os/signal"
	"sender-agent/discovery"
	metrics2 "sender-agent/metrics"
	"sender-agent/service"
	"sender-agent/subscriber"
	"strconv"
	"syscall"
	"time"
)

//BandwidthCounter tracks incoming and outgoing data transferred by the local peer.
var BandCounter *metrics.BandwidthCounter

func main() {

	err := godotenv.Load("sender.env")

	if err != nil {
		log.Println("Error loading sender.env file")
	}

	//set frequency of metrics sent from .env file, if an error occurs set this to 60s
	sendFrequency, err := strconv.Atoi(os.Getenv("SEND_FREQUENCY"))
	if err != nil {
		sendFrequency = 60
	}

	const (
		discoveryName = "discoveryRoom"
		roomPing      = "ping"
		roomSystem    = "system"
		roomCpu       = "cpu"
		roomRam       = "ram"
		roomTcp       = "tcp"
		roomBand      = "bandwidth"
	)

	context := context.Background()

	node := createHost()

	//return a new pubsub Service using the GossipSub router
	pubsub := subscriber.NewPubSubService(context, node)

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

	// setup local mDNS discovery
	discovery.SetupDiscovery(node, discoveryName)

	sender := service.NewSenderService(node, metrics2.LocalIP(), BandCounter, sendFrequency)

	//send Peer-System-Information
	go sender.SendPeerInfo(systemTopic, context, &discovery.PeerList)
	//send CPU information on a separated thread
	go sender.SendCpuInfo(cpuTopic, context, &discovery.PeerList)
	//send RAM information on a separated thread
	go sender.SendRamInfo(ramTopic, context, &discovery.PeerList)
	//send tcp status on a separated thread
	go sender.SendTCPstatus(tcpTopic, context, &discovery.PeerList)
	//Get Bandwidth between local Peer and other connected Peers on a separated thread
	go sender.GetBandWidthForActivePeer(systemSubscr, context, bandTopic)

	//Run the program till its stopped (forced)
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	fmt.Println("Received signal, shutting down...")
}

func createHost() host.Host {
	//return a tracker for the Bandwidth of the local Peer
	BandCounter = metrics.NewBandwidthCounter()
	// create a new libp2p Host that listens on a TCP port
	node, err := libp2p.New(libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/0"), libp2p.BandwidthReporter(BandCounter))
	//if an error appear we try again after 60 second
	if err != nil {
		time.Sleep(60 * time.Second)
		createHost()
	}
	return node
}
