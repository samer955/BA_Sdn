package main

import (
	"context"
	"database/sql"
	"db-agent/discovery"
	"db-agent/repository"
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

// BandwidthCounter tracks incoming and outgoing data transferred by the local peer.
var BandCounter *metrics.BandwidthCounter

func main() {

	const (
		discoveryTag = "discoveryRoom"
		roomSystem   = "system"
		roomCpu      = "cpu"
		roomRam      = "ram"
		roomPing     = "ping"
		roomTcp      = "tcp"
		host         = "localhost"
		port         = 5432
		user         = "user"
		password     = "password"
		dbname       = "database"
	)

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		panic(err)
	} else {
		fmt.Println("Successfully connected!")
	}

	//initialize Repository and DataColector
	repo := repository.NewPostGresRepository(db)
	repo.Migrate()
	collector := service.NewDataCollectorService(BandCounter, repo)
	context := context.Background()

	node := createHost()

	//create a new PubSub Service using the GossipSub router
	pubsub := subscriber.NewPubSubService(context, node)

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

	// setup local mDNS discovery
	discovery.SetupDiscovery(node, discoveryTag)

	//read System Information of peers in a separated thread
	go collector.ReadSystemInfo(systemSubscribe, context, node)
	//read cpu information of peers in a separated thread
	go collector.ReadCpuInformation(cpuSubscribe, context, node)
	//read ram information of peers in a separated thread
	go collector.ReadRamInformation(ramSubscribe, context, node)
	//read all the Ping Status from the other Peers
	go collector.ReadPingStatus(pingSubscribe, context, node)
	//read TCP infos from other Peers in a separated thread
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
