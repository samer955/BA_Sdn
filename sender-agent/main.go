package main

import (
	"context"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/metrics"
	"log"
	"net"
	"os"
	"os/signal"
	"sender-agent/discovery"
	"sender-agent/service"
	"sender-agent/subscriber"
	"syscall"
	"time"
)

//BandwidthCounter tracks incoming and outgoing data transferred by the local peer.
var BandCounter *metrics.BandwidthCounter

func main() {

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

	sender := service.NewSenderService(node, GetLocalIP(), BandCounter)

	//send Peer-System-Information
	go sender.SendPeerInfo(systemTopic, context, &discovery.PeerList)
	//send CPU information on a separated thread
	go sender.SendCpuInfo(cpuTopic, context, &discovery.PeerList)
	//send RAM information on a separated thread
	go sender.SendRamInfo(ramTopic, context, &discovery.PeerList)
	//send tcp status on a separated thread
	go sender.SendTCPstatus(tcpTopic, context, &discovery.PeerList)

	go sender.GetBandWidthForActivePeer(systemSubscr, context, bandTopic)

	//Run the program till its stopped (forced)
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	fmt.Println("Received signal, shutting down...")
}

func GetLocalIP() string {
	// testing with  198.18.0.0/15 , see https://www.iana.org/assignments/iana-ipv4-special-registry/iana-ipv4-special-registry.xhtml
	conn, err := net.Dial("udp", "198.18.0.30:80")
	if err != nil {
		log.Printf("Cannot use UDP: %s", err.Error())
		return "0.0.0.0"
	}

	defer conn.Close()

	if addr, ok := conn.LocalAddr().(*net.UDPAddr); ok {
		return addr.IP.String()
	}
	return ""
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
