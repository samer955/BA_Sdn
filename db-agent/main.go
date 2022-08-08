package main

import (
	"context"
	"db-agent/discovery"
	"db-agent/service"
	"db-agent/subscriber"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/metrics"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var BandCounter *metrics.BandwidthCounter

func main() {

	const (
		roomTime = "latency"
		roomCpu  = "cpu"
		roomRam  = "ram"
		roomPing = "ping"
		roomTcp  = "tcp"
	)

	context := context.Background()

	node := createHost()

	//return a new Pubsub Service using the GossipSub router
	_ = subscriber.NewPubSubService(context, node)

	pingTopic := subscriber.JoinTopic(roomPing)
	pingSubscribe := subscriber.Subscribe(pingTopic)

	timeTopic := subscriber.JoinTopic(roomTime)
	timeSubscribe := subscriber.Subscribe(timeTopic)

	tcpTopic := subscriber.JoinTopic(roomTcp)
	tcpSubscribe := subscriber.Subscribe(tcpTopic)

	cpuTopic := subscriber.JoinTopic(roomCpu)
	cpuSubscribe := subscriber.Subscribe(cpuTopic)

	ramTopic := subscriber.JoinTopic(roomRam)
	ramSubscribe := subscriber.Subscribe(ramTopic)

	// setup local mDNS discovery
	discovery.SetupDiscovery(node)

	collector := service.NewDataCollector()

	//read System Information of peers in a separated thread
	go collector.ReadSystemInfo(timeSubscribe, context, node)
	//read cpu information of peers in a separated thread
	go collector.ReadCpuInformation(cpuSubscribe, context, node)
	//read ram information of peers in a separated thread
	go collector.ReadRamInformation(ramSubscribe, context, node)
	//read all the Ping Status from the other Peers
	go collector.ReadPingStatus(pingSubscribe, context, node)
	//go service.ReadBandwidth(BandCounter, &PeerList)
	go collector.ReadTCPstatus(tcpSubscribe, context, node)

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
	//if an error appears we try again after 60 second
	if err != nil {
		time.Sleep(60 * time.Second)
		createHost()
	}
	return node
}
