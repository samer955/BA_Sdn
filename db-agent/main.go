package main

import (
	"context"
	"db-agent/receiver"
	"db-agent/subscriber"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/metrics"
	"github.com/libp2p/go-libp2p-core/peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const discoveryName = "discoveryRoom"

var PeerList []peer.AddrInfo
var PingTopic *pubsub.Topic
var BandCounter *metrics.BandwidthCounter

func main() {

	const roomTime = "latency"
	const roomCpu = "cpu"
	const roomRam = "ram"
	const roomPing = "ping"

	context := context.Background()
	node := createHost()

	//return a new Pubsub Service using the GossipSub router
	_ = subscriber.NewPubSubService(context, node)

	PingTopic = subscriber.JoinTopic(roomPing)
	pingSubscribe := subscriber.Subscribe(PingTopic)

	timeTopic := subscriber.JoinTopic(roomTime)
	timeSubscribe := subscriber.Subscribe(timeTopic)

	//cpuTopic := subscriber.JoinTopic(roomCpu)
	//cpuSubscribe := subscriber.Subscribe(cpuTopic)
	//
	//ramTopic := subscriber.JoinTopic(roomRam)
	//ramSubscribe := subscriber.Subscribe(ramTopic)

	// setup local mDNS discovery
	setupDiscovery(node)

	//read System Information of peers in a separated thread
	go receiver.ReadSystemInfo(timeSubscribe, context, node)
	//read cpu information of peers in a separated thread
	//go receiver.ReadCpuInformation(cpuSubscribe, context, node)
	//read ram information of peers in a separated thread
	//go receiver.ReadRamInformation(ramSubscribe, context, node)
	//read all the Ping Status from the other Peers
	go receiver.ReadPingStatus(pingSubscribe, context, node)
	go receiver.ReadBandwidth(BandCounter, &PeerList)

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

type discoveryNotifee struct {
	node host.Host
}

func (d *discoveryNotifee) HandlePeerFound(info peer.AddrInfo) {
	fmt.Printf("discovered a new peer %s\n", info.ID.Pretty())

	if d.node.ID().Pretty() != info.ID.Pretty() {
		d.node.Connect(context.Background(), info)
		PeerList = append(PeerList, info)

		log.Printf("connected to Peer %s ", info.ID.Pretty())
	}
}

func setupDiscovery(node host.Host) error {
	discovery := mdns.NewMdnsService(node, discoveryName, &discoveryNotifee{node: node})
	start := discovery.Start()

	//If any error is returned try again in 1min
	if start != nil {
		time.Sleep(60 * time.Second)
		setupDiscovery(node)
	}
	return start
}
