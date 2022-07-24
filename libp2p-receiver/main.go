package main

import (
	"context"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/libp2p/go-libp2p"
	_ "github.com/libp2p/go-libp2p-core/host"
	host2 "github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
	tls "github.com/libp2p/go-libp2p/p2p/security/tls"
	"github.com/libp2p/go-libp2p/p2p/transport/tcp"
	"github.com/libp2p/go-libp2p/p2p/transport/websocket"
	"libp2p-receiver/receiver"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const discoveryName = "discoveryRoom"
const discoveryInterval = time.Minute

var PeerList []peer.AddrInfo

func main() {

	const roomTime = "latency"
	const roomCpu = "cpu"
	const roomRam = "ram"

	context := context.Background()

	transports := libp2p.ChainOptions(
		libp2p.Transport(tcp.NewTCPTransport),
		libp2p.Transport(websocket.New),
	)

	security := libp2p.Security(tls.ID, tls.New)

	listenAddrs := libp2p.ListenAddrStrings(
		"/ip4/0.0.0.0/tcp/0",
		"/ip4/0.0.0.0/tcp/0/ws",
	)

	// create a new libp2p Host that listens on a random TCP port
	node, err := libp2p.New(transports, listenAddrs, security, libp2p.Ping(false))
	if err != nil {
		panic(err)
	}

	//create a new pubsub Service using the GossipSub router
	ps, err := pubsub.NewGossipSub(context, node)
	if err != nil {
		panic(err)
	}

	timeTopic, err := ps.Join(roomTime)

	if err != nil {
		log.Println("Error while subscribing in the TIME-Topic")
	} else {
		log.Println("Subscribed on", roomTime)
		log.Println("topicID", timeTopic.String())
	}

	cpuTopic, err := ps.Join(roomCpu)

	if err != nil {
		log.Println("Error while subscribing in the CPU-Topic")
	} else {
		log.Println("Subscribed on", roomCpu)
		log.Println("topicID", cpuTopic.String())
	}

	ramTopic, err := ps.Join(roomRam)

	if err != nil {
		log.Println("Error while subscribing in the CPU-Topic")
	} else {
		log.Println("Subscribed on", roomRam)
		log.Println("topicID", ramTopic.String())
	}

	// setup local mDNS discovery
	if err := setupDiscovery(node); err != nil {
		panic(err)
	}

	subscribe, err := timeTopic.Subscribe()
	if (err) != nil {
		log.Println("cannot subscribe to: ", timeTopic.String())
	}
	subscribe2, err := cpuTopic.Subscribe()
	if (err) != nil {
		log.Println("cannot subscribe to: ", cpuTopic.String())
	}
	subscribe3, err := ramTopic.Subscribe()
	if (err) != nil {
		log.Println("cannot subscribe to: ", ramTopic.String())
	}

	//read timestamp of peers in a separated thread
	go receiver.ReadTimeMessages(subscribe, context, node)
	//read cpu information of peers in a separated thread
	go receiver.ReadCpuInformation(subscribe2, context, node)
	//read ram information of peers in a separated thread
	go receiver.ReadRamInformation(subscribe3, context, node)

	//Run the program till its stopped
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	fmt.Println("Received signal, shutting down...")
}

type discoveryNotifee struct {
	node host2.Host
}

func (d *discoveryNotifee) HandlePeerFound(info peer.AddrInfo) {
	fmt.Printf("discovered a new peer %s\n", info.ID.Pretty())

	if d.node.ID().Pretty() != info.ID.Pretty() {
		d.node.Connect(context.Background(), info)
		PeerList = append(PeerList, info)

		log.Printf("connected to Peer %s ", info.ID.Pretty())
		go receiver.SendPing(context.Background(), d.node, info)
	}
}

func setupDiscovery(node host2.Host) error {
	discovery := mdns.NewMdnsService(node, discoveryName, &discoveryNotifee{node: node})
	return discovery.Start()
}
