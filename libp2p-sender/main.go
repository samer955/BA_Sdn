package main

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	"libp2p-sender/components"
	"libp2p-sender/discovery"
	"libp2p-sender/service"
	"libp2p-sender/subscriber"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	const (
		discoveryName = "discoveryRoom"
		roomPing      = "ping"
		roomTime      = "latency"
		roomCpu       = "cpu"
		roomRam       = "ram"
		roomTcp       = "tcp"
	)

	context := context.Background()

	node := createHost()

	//return a new pubsub Service using the GossipSub router
	_ = subscriber.NewPubSubService(context, node)

	timeTopic := subscriber.JoinTopic(roomTime)
	_ = subscriber.Subscribe(timeTopic)

	pingTopic := subscriber.JoinTopic(roomPing)
	_ = subscriber.Subscribe(pingTopic)

	discovery.SetPingTopic(pingTopic)

	cpuTopic := subscriber.JoinTopic(roomCpu)
	_ = subscriber.Subscribe(cpuTopic)

	ramTopic := subscriber.JoinTopic(roomRam)
	_ = subscriber.Subscribe(ramTopic)

	tcpTopic := subscriber.JoinTopic(roomTcp)
	_ = subscriber.Subscribe(tcpTopic)

	// setup local mDNS discovery
	discovery.SetupDiscovery(node, discoveryName)

	ipAddress := GetLocalIP()

	peer_sys := components.NewPeerInfo(ipAddress, node.ID().Pretty(), os.Getenv("ROLE_HOST"))
	peer_cpu := components.NewCpu(ipAddress, node.ID().Pretty())
	peer_ram := components.NewRam(ipAddress, node.ID().Pretty())
	peer_tcp := components.NewTCPstatus(ipAddress)

	//send timestamp on a separated thread
	go service.SendPeerInfo(timeTopic, context, peer_sys, &discovery.PeerList)
	//send CPU information on a separated thread
	go service.SendCpuInfo(cpuTopic, context, peer_cpu, &discovery.PeerList)
	//send RAM information on a separated thread
	go service.SendRamInfo(ramTopic, context, peer_ram, &discovery.PeerList)
	//send tcp status on a separated thread
	go service.SendTCPstatus(tcpTopic, context, peer_tcp, &discovery.PeerList)

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
	// create a new libp2p Host that listens on a TCP port
	node, err := libp2p.New(libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/0"))
	//if an error appear we try again after 60 second
	if err != nil {
		time.Sleep(60 * time.Second)
		createHost()
	}
	return node
}
