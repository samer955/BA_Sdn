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

func main() {

	const roomTime = "latency"
	const roomCpu = "cpu"
	const roomRam = "ram"
	const roomPing = "ping"

	context := context.Background()

	//return a tracker for the Bandwidth of the local Peer
	//counter := new(metrics.BandwidthCounter)

	node := createHost()

	//return a new Pubsub Service using the GossipSub router
	_ = subscriber.NewPubSubService(context, node)

	PingTopic = subscriber.JoinTopic(roomPing)
	pingSubscribe := subscriber.Subscribe(PingTopic)

	timeTopic := subscriber.JoinTopic(roomTime)
	timeSubscribe := subscriber.Subscribe(timeTopic)

	cpuTopic := subscriber.JoinTopic(roomCpu)
	cpuSubscribe := subscriber.Subscribe(cpuTopic)

	ramTopic := subscriber.JoinTopic(roomRam)
	ramSubscribe := subscriber.Subscribe(ramTopic)

	// setup local mDNS discovery
	if err := setupDiscovery(node); err != nil {
		time.Sleep(60 * time.Second)
		setupDiscovery(node)
	}
	//read Timestamp of peers in a separated thread
	go receiver.ReadTimeMessages(timeSubscribe, context, node)
	//read cpu information of peers in a separated thread
	go receiver.ReadCpuInformation(cpuSubscribe, context, node)
	//read ram information of peers in a separated thread
	go receiver.ReadRamInformation(ramSubscribe, context, node)
	//read all the Ping Status from the other Peers
	go receiver.ReadPingStatus(pingSubscribe, context, node)

	//go readBandwidth(counter)

	//Run the program till its stopped (forced)
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	fmt.Println("Received signal, shutting down...")
}

func createHost() host.Host {
	// create a new libp2p Host that listens on a TCP port
	node, err := libp2p.New(libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/0"))
	//if an error appear we try again after 60 second
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
		//go receiver.SendPing(context.Background(), d.node, info)
	}
}

func setupDiscovery(node host.Host) error {
	discovery := mdns.NewMdnsService(node, discoveryName, &discoveryNotifee{node: node})
	return discovery.Start()
}

func readBandwidth(counter *metrics.BandwidthCounter) {
	for {
		total := counter.GetBandwidthTotals()
		fmt.Println(total)
		mapPeer := counter.GetBandwidthByPeer()
		fmt.Println(mapPeer)
		time.Sleep(60 * time.Second)
	}
}
