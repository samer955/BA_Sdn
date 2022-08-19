package service

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/p2p/protocol/ping"
	"os"
	"sender-agent/metrics"
	"sender-agent/subscriber"
	"time"
)

type Sender struct {
	node host.Host
	ip   string
}

func NewSenderService(node host.Host, ip string) *Sender {
	return &Sender{
		node: node,
		ip:   ip}
}

//SendPeerInfo will send system Information of this Peer and periodically a timestamp
//in order to calculate the latency in ms between service and service
func (s *Sender) SendPeerInfo(topic *pubsub.Topic, context context.Context, list *[]peer.AddrInfo) {
	peer_sys := metrics.NewPeerInfo(s.ip, s.node.ID(), os.Getenv("ROLE_HOST"))
	for {
		if len(*list) == 0 {
			continue
		}
		sendPeerInfo(topic, context, peer_sys)

		//wait 20 seconds before send another timestamp
		time.Sleep(20 * time.Second)
	}
}

func sendPeerInfo(topic *pubsub.Topic, context context.Context, peer *metrics.PeerInfo) {
	fmt.Println("sending time")

	//Set the time when the message is sent
	peer.UUID = uuid.New().String()
	peer.OnlineUser = metrics.GetNumberOfOnlineUser()
	peer.Time = metrics.TimeFromServer()

	err := subscriber.Publish(peer, context, topic)
	if err != nil {
		fmt.Println("Error publishing content ", err.Error())
	}
}

// SendCpuInfo function will send periodically information about the CPU
func (s *Sender) SendCpuInfo(topic *pubsub.Topic, context context.Context, peers *[]peer.AddrInfo) {
	cpu := metrics.NewCpu(s.ip, s.node.ID().Pretty())
	for {
		if len(*peers) == 0 {
			continue
		}
		sendCpuInfo(topic, context, cpu)

		time.Sleep(15 * time.Second)
	}
}

func sendCpuInfo(topic *pubsub.Topic, context context.Context, cpu *metrics.Cpu) {
	fmt.Println("sending cpu")

	cpu.UUID = uuid.New().String()
	//Update every 10s CPU Usages in %
	cpu.UpdateCpuPercentage()

	//publish the cpu data
	err := subscriber.Publish(cpu, context, topic)

	if err != nil {
		fmt.Println("Error publishing content ", err.Error())
	}
}

// SendRamInfo function will send periodically information about the actual RAM Percentage
func (s *Sender) SendRamInfo(topic *pubsub.Topic, context context.Context, peers *[]peer.AddrInfo) {
	ram := metrics.NewRam(s.ip, s.node.ID().Pretty())
	for {
		if len(*peers) == 0 {
			continue
		}
		sendRamInfo(topic, context, ram)
		time.Sleep(5 * time.Second)
	}
}

func sendRamInfo(topic *pubsub.Topic, context context.Context, ram *metrics.Ram) {
	fmt.Println("sending ram")

	ram.UUID = uuid.New().String()
	//Update every 10s RAM usages in %
	ram.UpdateRamPercentage()

	//publish the ram data
	err := subscriber.Publish(ram, context, topic)
	if err != nil {
		fmt.Println("Error publishing content ", err.Error())
	}
}

//SendPing function send a Ping every 60s to the discovered Peer
//This function is used in order to reach the others nodes in an active way getting
//a bool if the ping was successfully (true if a node is reachable, false if not) and an RTT in ms
func SendPing(ctx context.Context, host host.Host, target peer.AddrInfo, topic *pubsub.Topic) {

	var pingDeadline = 10

	status := metrics.NewPingStatus(host.ID().Pretty(), target.ID.Pretty())
	//The Ping function return a channel that still open till the context is alive
	ch := ping.Ping(ctx, host, target.ID)

	for {
		//after 10 negative Ping stop to ping the Peer
		if pingDeadline == 0 {
			fmt.Printf("Stopped ping from %s to %s\n", status.Source, status.Target)
			return
		}
		res := <-ch

		status.SetPingStatus(res, &pingDeadline)

		//publish the status of the Ping in the topic
		subscriber.Publish(status, ctx, topic)

		//Next Ping in 1 Min
		time.Sleep(5 * time.Second)
	}
}

func (s *Sender) SendTCPstatus(topic *pubsub.Topic, context context.Context, peers *[]peer.AddrInfo) {
	tcp := metrics.NewTCPstatus(s.ip)
	for {
		if len(*peers) == 0 {
			continue
		}
		sendTCPstatus(topic, context, tcp)
		time.Sleep(15 * time.Second)
	}
}

func sendTCPstatus(topic *pubsub.Topic, context context.Context, tcpIfo *metrics.TCPstatus) {
	fmt.Println("sending tcp Queue size")

	tcpIfo.UUID = uuid.New().String()
	tcpIfo.QueueSize = metrics.TcpQueueSize()
	received, sent := metrics.TcpSegmentsNumber()

	tcpIfo.Received = received
	tcpIfo.Sent = sent
	tcpIfo.Time = metrics.TimeFromServer()

	err := subscriber.Publish(tcpIfo, context, topic)
	if err != nil {
		fmt.Println("Error publishing content ", err.Error())
	}
}
