package main

import (
	"context"
	"database/sql"
	"db-agent/discovery"
	"db-agent/repository"
	"db-agent/service"
	"db-agent/subscriber"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	const (
		discoveryTag = "discoveryRoom"
		roomSystem   = "system"
		roomCpu      = "cpu"
		roomRam      = "ram"
		roomPing     = "ping"
		roomTcp      = "tcp"
		roomBand     = "bandwidth"
	)

	err := godotenv.Load("db.env")

	if err != nil {
		log.Println("Error loading db.env file")
	}

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))

	db, err := sql.Open(os.Getenv("DB_DRIVER"), psqlInfo)

	if err != nil {
		panic(err)
	} else {
		fmt.Println("Successfully connected!")
	}

	node := createHost()
	context := context.Background()

	//initialize Repository and DataCollector
	repo := repository.NewPostGresRepository(db)
	repo.Migrate()
	receiver := service.NewDataCollectorService(node, repo)

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

	bandTopic := pubsub.JoinTopic(roomBand)
	bandSubscribe := pubsub.Subscribe(bandTopic)

	// setup local mDNS discovery
	discovery.SetupDiscovery(node, discoveryTag)

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
	//read Bandwidth infos from other Peers in a separated thread
	go receiver.ReadBandwidth(bandSubscribe, context)

	//Run the program till its stopped (forced)
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	fmt.Println("Received signal, shutting down...")
}

func createHost() host.Host {
	// create a new libp2p Host that listens on a TCP port
	node, err := libp2p.New(libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/0"))
	//if an error appears we try again after 60 second
	if err != nil {
		time.Sleep(60 * time.Second)
		createHost()
	}
	return node
}
