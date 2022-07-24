package main

import (
	"context"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
	"libp2p-sender/sender"
	"libp2p-sender/subscriber"
	"libp2p-sender/variables"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

const discoveryName = "discoveryRoom"

var PeerList []peer.AddrInfo
var PingTopic *pubsub.Topic

func main() {

	const roomPing = "ping"
	const roomTime = "latency"
	const roomCpu = "cpu"
	const roomRam = "ram"

	context := context.Background()

	// create a new libp2p Host that listens on a random TCP port
	node, err := libp2p.New(libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/0"))
	if err != nil {
		panic(err)
	}

	// setup local mDNS discovery
	if err := setupDiscovery(node); err != nil {
		panic(err)
	}

	//return a new pubsub Service using the GossipSub router
	_ = subscriber.NewPubSubService(context, node)

	PingTopic = subscriber.JoinTopic(roomPing)
	_ = subscriber.Subscribe(PingTopic)

	timeTopic := subscriber.JoinTopic(roomTime)
	_ = subscriber.Subscribe(timeTopic)

	cpuTopic := subscriber.JoinTopic(roomCpu)
	_ = subscriber.Subscribe(cpuTopic)

	ramTopic := subscriber.JoinTopic(roomRam)
	_ = subscriber.Subscribe(ramTopic)

	peer_lat := variables.NewPeerInfo(GetLocalIP(), node.ID().Pretty())
	peer_cpu := variables.NewCpu(GetLocalIP(), node.ID().Pretty())
	peer_ram := variables.NewRam(GetLocalIP(), node.ID().Pretty())

	//send timestamp on a separated thread
	go sender.SendTimeMessage(timeTopic, context, peer_lat, PeerList)
	//send CPU information on a separated thread
	go sender.SendCpuInformation(cpuTopic, context, peer_cpu, PeerList)
	//send RAM information on a separated thread
	go sender.SendRamInformation(ramTopic, context, peer_ram, PeerList)

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

type discoveryNotifee struct {
	node host.Host
}

func (d *discoveryNotifee) HandlePeerFound(info peer.AddrInfo) {
	fmt.Printf("discovered a new peer %s\n", info.ID.Pretty())

	if d.node.ID().Pretty() != info.ID.Pretty() {
		d.node.Connect(context.Background(), info)
		PeerList = append(PeerList, info)

		log.Printf("connected to Peer %s ", info.ID.Pretty())
		sender.SendPing(context.Background(), d.node, info, PingTopic)
	}
}

func setupDiscovery(node host.Host) error {
	discovery := mdns.NewMdnsService(node, discoveryName, &discoveryNotifee{node: node})
	return discovery.Start()
}
