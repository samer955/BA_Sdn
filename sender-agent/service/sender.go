package service

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/p2p/protocol/ping"
	"log"
	"sender-agent/metrics"
	"sender-agent/node"
	"sender-agent/subscriber"
	"time"
)

type Sender struct {
	Node      node.Node
	Frequency int
}

//SendPeerInfo will send system Information of this Peer and periodically a timestamp
//in order to calculate the latency in ms between service and service
func (s *Sender) SendPeerInfo(topic *pubsub.Topic, context context.Context) {
	systemInfo := metrics.NewSystemInfo(s.Node.Ip, s.Node.Host.ID(), s.Node.Role, s.Node.Network)

	for {
		//the key of the local node is also present in the peerstore, so we check if this is == 1
		//the node waits to connect to another node before sending his information
		if len(s.Node.Host.Peerstore().Peers()) == 1 {
			continue
		}
		sendPeerInfo(topic, context, systemInfo)

		//wait some time defined in the frequency before send another systeminfo
		time.Sleep(time.Duration(s.Frequency) * time.Second)
	}
}

func sendPeerInfo(topic *pubsub.Topic, context context.Context, systemInfo *metrics.SystemInfo) {

	systemInfo.UpdateLoggedInUser()
	systemInfo.UUID = uuid.New().String()
	systemInfo.Time = time.Now()

	err := subscriber.Publish(systemInfo, context, topic)
	if err != nil {
		log.Println("Error publishing content ", err.Error())
		return
	}
	log.Println("sending system info...")
}

// SendCpuInfo function will send periodically information about the CPU
func (s *Sender) SendCpuInfo(topic *pubsub.Topic, context context.Context) {
	cpu := metrics.NewCpu(s.Node.Ip, s.Node.Host.ID().Pretty())
	for {
		if len(s.Node.Host.Peerstore().Peers()) == 1 {
			continue
		}
		sendCpuInfo(topic, context, cpu)

		time.Sleep(time.Duration(s.Frequency) * time.Second)
	}
}

func sendCpuInfo(topic *pubsub.Topic, context context.Context, cpu *metrics.Cpu) {

	cpu.UUID = uuid.New().String()
	//Update CPU % Usages
	cpu.UpdateCpuPercentage()

	//publish the cpu data
	err := subscriber.Publish(cpu, context, topic)

	if err != nil {
		log.Println("Error publishing content ", err.Error())
		return
	}
	log.Println("sending cpu...")
}

// SendRamInfo function will send periodically information about the actual RAM Percentage
func (s *Sender) SendRamInfo(topic *pubsub.Topic, context context.Context) {
	ram := metrics.NewRam(s.Node.Ip, s.Node.Host.ID().Pretty())
	for {
		if len(s.Node.Host.Peerstore().Peers()) == 1 {
			continue
		}
		sendRamInfo(topic, context, ram)
		time.Sleep(time.Duration(s.Frequency) * time.Second)
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
		return
	}
	log.Println("sending ram...")
}

//Ping function send a Ping every 45s to the new discovered Peer
//This function is used in order to reach the others nodes in an active way getting
//a bool if the ping was successfully (true if a Node is reachable, false if not) and an RTT in ms
func Ping(ctx context.Context, host host.Host, target peer.AddrInfo, topic *pubsub.Topic) {

	var pingDeadlineLimit = 10
	var actualNegativePing = 0

	status := metrics.NewPingStatus(host.ID().Pretty(), target.ID.Pretty())
	//The Ping function return a channel that still open till the context is alive
	ch := ping.Ping(ctx, host, target.ID)

	for {
		//after 10 negative Ping stop to ping the Peer
		if actualNegativePing == pingDeadlineLimit {
			log.Printf("Stopped ping from %s to %s\n", status.Source, status.Target)
			return
		}
		res := <-ch

		status.SetPingStatus(res, &actualNegativePing)

		//publish the status of the Ping in the topic
		err := subscriber.Publish(status, ctx, topic)
		if err != nil {
			log.Println("Error publishing content ", err.Error())
		}
		//Next Ping in 45s
		time.Sleep(45 * time.Second)
	}
}

func (s *Sender) SendTCPstatus(topic *pubsub.Topic, context context.Context) {
	tcpStatus := metrics.NewTCPstatus(s.Node.Ip)
	for {
		if len(s.Node.Host.Peerstore().Peers()) == 1 {
			continue
		}
		sendTCPstatus(topic, context, tcpStatus)
		time.Sleep(time.Duration(s.Frequency) * time.Second)
	}
}

func sendTCPstatus(topic *pubsub.Topic, context context.Context, tcpStatus *metrics.TCPstatus) {

	tcpStatus.UUID = uuid.New().String()
	tcpStatus.TcpQueueSize()
	tcpStatus.TcpSegmentsNumber()
	tcpStatus.Time = time.Now()

	err := subscriber.Publish(tcpStatus, context, topic)
	if err != nil {
		log.Println("Error publishing content ", err.Error())
		return
	}
	log.Println("sending TCP-info...")
}

//GetBandWidthForActivePeer listens on the sytemtopic to get the information about an online Peer to get
//the Bandwidth between them
func (s *Sender) GetBandWidthForActivePeer(subscribe *pubsub.Subscription, context context.Context, topic *pubsub.Topic) {
	for {
		message, err := subscribe.Next(context)
		if err != nil {
			log.Println("cannot read from topic")
		} else {
			if message.ReceivedFrom.String() != s.Node.Host.ID().Pretty() {
				targetPeer := new(metrics.SystemInfo)
				json.Unmarshal(message.Data, targetPeer)
				s.getBandwidth(targetPeer, topic, context)
			}
		}
	}
}

func (s *Sender) getBandwidth(peer *metrics.SystemInfo, topic *pubsub.Topic, ctx context.Context) {

	bandwidth := metrics.NewBandWidth(s.Node.Ip, s.Node.Host.ID().Pretty())

	resultBand := s.Node.Bandcounter.GetBandwidthForPeer(peer.Id)

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
		return
	}
	log.Println("sending Bandwidth...")
}
