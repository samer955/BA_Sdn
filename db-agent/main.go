package main

import (
	"context"
	"db-agent/config"
	"db-agent/discovery"
	"db-agent/node"
	"db-agent/repository"
	"db-agent/service"
	"db-agent/subscriber"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	const (
		discoveryName = "discoveryRoom"
		roomSystem    = "system"
		roomCpu       = "cpu"
		roomRam       = "ram"
		roomPing      = "ping"
		roomTcp       = "tcp"
		roomBand      = "bandwidth"
	)

	var (
		node node.Node
		mdns error
	)

	context := context.Background()

	config.LoadEnv()
	config := config.GetConfig()

	//initialize node
	node.StartNode()

	//initialize Repository and Receiver service
	repo := repository.NewPostGresRepository(config.Connection)
	repo.Migrate()
	receiver := service.Receiver{Repository: repo, Node: node}

	//create a new PubSub Service using the GossipSub router
	pubsub := subscriber.NewPubSubService(context, node)

	//subscribe to all the topics
	pingTopic := pubsub.JoinTopic(roomPing)
	pingSubscribe := pubsub.Subscribe(pingTopic)

	systemTopic := pubsub.JoinTopic(roomSystem)
	systemSubscribe := pubsub.Subscribe(systemTopic)

	tcpTopic := pubsub.JoinTopic(roomTcp)
	tcpSubscribe := pubsub.Subscribe(tcpTopic)

	cpuTopic := pubsub.JoinTopic(roomCpu)
	cpuSubscribe := pubsub.Subscribe(cpuTopic)

	ramTopic := pubsub.JoinTopic(roomRam)
	ramSubscribe := pubsub.Subscribe(ramTopic)

	bandTopic := pubsub.JoinTopic(roomBand)
	bandSubscribe := pubsub.Subscribe(bandTopic)

	//setup local mDNS discovery. If an error occurs try again in 60 seconds
	for {
		mdns = discovery.SetupDiscovery(node, discoveryName)
		if mdns == nil {
			break
		}
		log.Println("unable to start the mDNS discovery, next try in 60 seconds ...")
		time.Sleep(60 * time.Second)
	}

	//read System Information of peers in a separated thread
	go receiver.ReadSystemInfo(systemSubscribe, context)
	//read cpu information of peers in a separated thread
	go receiver.ReadCpuInformation(cpuSubscribe, context)
	//read ram information of peers in a separated thread
	go receiver.ReadRamInformation(ramSubscribe, context)
	//read all the Ping Status from the other Peers
	go receiver.ReadPingStatus(pingSubscribe, context)
	//read TCP infos from other Peers in a separated thread
	go receiver.ReadTCPstatus(tcpSubscribe, context)
	//read Bandwidth from other Peers in a separated thread
	go receiver.ReadBandwidth(bandSubscribe, context)

	//Run the program till its stopped (forced)
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	fmt.Println("Received signal, shutting down...")
}
