package main

import (
	"context"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/libp2p/go-libp2p"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"libp2p-sender/discovery"
	"libp2p-sender/sender"
	"libp2p-sender/variables"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const discoveryInterval = time.Minute

func main() {

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
	if err := discovery.SetupDiscovery(node); err != nil {
		panic(err)
	}

	//create a new pubsub Service using the GossipSub router
	ps, err := pubsub.NewGossipSub(context, node)
	if err != nil {
		panic(err)
	}

	timeTopic, err := ps.Join(roomTime)

	if err != nil {
		log.Println("Error while subscribing in the Time-Topic")
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
		log.Println("Error while subscribing in the RAM-Topic")
	} else {
		log.Println("Subscribed on", roomRam)
		log.Println("topicID", ramTopic.String())
	}

	subscribe, err := timeTopic.Subscribe()

	if (err) != nil {
		log.Println("cannot subscribe to: ", timeTopic.String())
	} else {
		log.Println("Subscribed to, " + subscribe.Topic())
	}

	subscribe2, err := cpuTopic.Subscribe()

	if (err) != nil {
		log.Println("cannot subscribe to: ", cpuTopic.String())
	} else {
		log.Println("Subscribed to, " + subscribe2.Topic())
	}

	//	peer_lat := variables.NewPeerInfo(GetLocalIP(), node.ID().Pretty())
	peer_cpu := variables.NewCpu(GetLocalIP(), node.ID().Pretty())
	//	peer_ram := variables.NewRam(GetLocalIP(), node.ID().Pretty())

	//send timestamp on a separated thread
	//	go sender.SendTimeMessage(timeTopic, context, peer_lat)
	//send CPU information on a separated thread
	go sender.SendCpuInformation(cpuTopic, context, peer_cpu)
	//send RAM information on a separated thread
	//	go sender.SendRamInformation(ramTopic, context, peer_ram)

	//Run the program till its stopped
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
