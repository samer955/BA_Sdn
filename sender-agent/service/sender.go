package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/libp2p/go-libp2p-core/host"
	metrics2 "github.com/libp2p/go-libp2p-core/metrics"
	"github.com/libp2p/go-libp2p-core/peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/p2p/protocol/ping"
	"log"
	"os"
	"sender-agent/metrics"
	"sender-agent/subscriber"
	"time"
)

type Sender struct {
	node    host.Host
	ip      string
	counter *metrics2.BandwidthCounter
}

func NewSenderService(node host.Host, ip string, counter *metrics2.BandwidthCounter) *Sender {
	return &Sender{
		node:    node,
		ip:      ip,
		counter: counter}
}

//SendPeerInfo will send system Information of this Peer and periodically a timestamp
//in order to calculate the latency in ms between service and service
func (s *Sender) SendPeerInfo(topic *pubsub.Topic, context context.Context, list *[]peer.AddrInfo) {
	peerSys := metrics.NewPeerInfo(s.ip, s.node.ID(), os.Getenv("ROLE_HOST"))
	for {
		if len(*list) == 0 {
			continue
		}
		sendPeerInfo(topic, context, peerSys)

		//wait 30 seconds before send another systeminfo
		time.Sleep(30 * time.Second)
	}
}

func sendPeerInfo(topic *pubsub.Topic, context context.Context, peer *metrics.PeerInfo) {

	peer.UUID = uuid.New().String()
	peer.OnlineUser = metrics.GetNumberOfOnlineUser()
	peer.Time = metrics.TimeFromServer()

	err := subscriber.Publish(peer, context, topic)
	if err != nil {
		fmt.Println("Error publishing content ", err.Error())
	}
	log.Println("sending system info...")
}

// SendCpuInfo function will send periodically information about the CPU
func (s *Sender) SendCpuInfo(topic *pubsub.Topic, context context.Context, peers *[]peer.AddrInfo) {
	cpu := metrics.NewCpu(s.ip, s.node.ID().Pretty())
	for {
		if len(*peers) == 0 {
			continue
		}
		sendCpuInfo(topic, context, cpu)

		time.Sleep(30 * time.Second)
	}
}

func sendCpuInfo(topic *pubsub.Topic, context context.Context, cpu *metrics.Cpu) {

	cpu.UUID = uuid.New().String()
	//Update CPU % Usages
	cpu.UpdateCpuPercentage()

	//publish the cpu data
	err := subscriber.Publish(cpu, context, topic)

	if err != nil {
		fmt.Println("Error publishing content ", err.Error())
	}
	log.Println("sending cpu...")
}

// SendRamInfo function will send periodically information about the actual RAM Percentage
func (s *Sender) SendRamInfo(topic *pubsub.Topic, context context.Context, peers *[]peer.AddrInfo) {
	ram := metrics.NewRam(s.ip, s.node.ID().Pretty())
	for {
		if len(*peers) == 0 {
			continue
		}
		sendRamInfo(topic, context, ram)
		time.Sleep(30 * time.Second)
	}
}

func sendRamInfo(topic *pubsub.Topic, context context.Context, ram *metrics.Ram) {

	ram.UUID = uuid.New().String()
	//Update RAM % Usage
	ram.UpdateRamPercentage()

	//publish the ram data
	err := subscriber.Publish(ram, context, topic)
	if err != nil {
		log.Println("Error publishing content ", err.Error())
	}
	log.Println("sending ram...")
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
			log.Printf("Stopped ping from %s to %s\n", status.Source, status.Target)
			return
		}
		res := <-ch

		status.SetPingStatus(res, &pingDeadline)

		//publish the status of the Ping in the topic
		subscriber.Publish(status, ctx, topic)

		//Next Ping in 45s
		time.Sleep(45 * time.Second)
	}
}

func (s *Sender) SendTCPstatus(topic *pubsub.Topic, context context.Context, peers *[]peer.AddrInfo) {
	tcp := metrics.NewTCPstatus(s.ip)
	for {
		if len(*peers) == 0 {
			continue
		}
		sendTCPstatus(topic, context, tcp)
		time.Sleep(30 * time.Second)
	}
}

func sendTCPstatus(topic *pubsub.Topic, context context.Context, tcpIfo *metrics.TCPstatus) {

	tcpIfo.UUID = uuid.New().String()
	tcpIfo.QueueSize = metrics.TcpQueueSize()
	received, sent := metrics.TcpSegmentsNumber()

	tcpIfo.Received = received
	tcpIfo.Sent = sent
	tcpIfo.Time = metrics.TimeFromServer()

	err := subscriber.Publish(tcpIfo, context, topic)
	if err != nil {
		log.Println("Error publishing content ", err.Error())
	}
	log.Println("sending TCP-info...")
}

// GetBandWidthForActivePeer listens on the sytemtopic to get the information about an online Peer in order to calculate
//the Bandwidth between them
func (s *Sender) GetBandWidthForActivePeer(subscribe *pubsub.Subscription, context context.Context, topic *pubsub.Topic) {
	for {
		message, err := subscribe.Next(context)
		if err != nil {
			log.Println("cannot read from topic")
		} else {
			if message.ReceivedFrom.String() != s.node.ID().Pretty() {
				peer := new(metrics.PeerInfo)
				json.Unmarshal(message.Data, peer)
				s.getBandwidth(peer, topic, context)
			}
		}
	}
}

func (s *Sender) getBandwidth(peer *metrics.PeerInfo, topic *pubsub.Topic, ctx context.Context) {

	bandwidth := metrics.NewBandWidth(s.ip, s.node.ID().Pretty())

	resultBand := s.counter.GetBandwidthForPeer(peer.Id)

	bandwidth.Target = peer.Ip
	bandwidth.TotalIn = resultBand.TotalIn
	bandwidth.TotalOut = resultBand.TotalOut
	bandwidth.RateIn = int(resultBand.RateIn)
	bandwidth.RateOut = int(resultBand.RateOut)
	bandwidth.Time = peer.Time
	bandwidth.UUID = uuid.New().String()

	err := subscriber.Publish(bandwidth, ctx, topic)
	if err != nil {
		log.Println("Error publishing content ", err.Error())
	}
	log.Println("sending Bandwidth...")
}
