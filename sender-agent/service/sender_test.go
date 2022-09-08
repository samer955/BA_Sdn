package service

import (
	"context"
	"encoding/json"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
	"github.com/stretchr/testify/assert"
	"sender-agent/config"
	"sender-agent/metrics"
	node2 "sender-agent/node"
	"sender-agent/subscriber"
	"testing"
)

func setupEnvironment(t *testing.T) (*pubsub.Topic, *pubsub.Subscription, context.Context, node2.Node, config.Config) {
	const roomTest = "test"
	var node node2.Node
	node.StartNode()
	ctx := context.Background()
	psub := subscriber.NewPubSubService(ctx, node)
	testTopic := psub.JoinTopic(roomTest)
	subsc := psub.Subscribe(testTopic)
	conf := config.Config{Frequency: 30, Role: "SENDER", Network: "HOME"}

	t.Cleanup(func() {
		node.Host.Close()
		ctx.Done()
		testTopic.Close()
		subsc.Cancel()
		node.Bandcounter.Reset()
	})
	return testTopic, subsc, ctx, node, conf
}

//Here it is tested if the system-information are sent/published
func TestSendPeerInfo(t *testing.T) {
	topic, subscr, ctx, _, _ := setupEnvironment(t)
	peerInfo := metrics.NewSystemInfo("1.1.1.1", "test_ID", "TEST", "HOME")

	sendPeerInfo(topic, ctx, peerInfo)
	message, _ := subscr.Next(ctx)
	peerToBytesConverted, _ := json.Marshal(peerInfo)

	assert.Equal(t, message.Data, peerToBytesConverted)
}

//Here it is tested if the cpu is sent/published
func TestSendCpuInfo(t *testing.T) {
	topic, subscr, ctx, _, _ := setupEnvironment(t)
	cpu := metrics.NewCpu("1.1.1.1", "test_ID")

	sendCpuInfo(topic, ctx, cpu)
	message, _ := subscr.Next(ctx)
	cpuToBytesConverted, _ := json.Marshal(cpu)

	assert.Equal(t, message.Data, cpuToBytesConverted)
}

//Here it is tested if the ram is sent/published
func TestSendRamInfo(t *testing.T) {
	topic, subscr, ctx, _, _ := setupEnvironment(t)
	ram := metrics.NewRam("1.1.1.1", "test_ID")

	sendRamInfo(topic, ctx, ram)
	message, _ := subscr.Next(ctx)
	ramToBytesConverted, _ := json.Marshal(ram)

	assert.Equal(t, message.Data, ramToBytesConverted)
}

//Here it is tested if the tcpStatus is sent/published
func TestSendTCPstatus(t *testing.T) {
	topic, subscr, ctx, _, _ := setupEnvironment(t)
	tcpStatus := metrics.NewTCPstatus("1.1.1.1")

	sendTCPstatus(topic, ctx, tcpStatus)
	message, _ := subscr.Next(ctx)
	ramToBytesConverted, _ := json.Marshal(tcpStatus)

	assert.Equal(t, message.Data, ramToBytesConverted)
}

type discoveryNotifee struct {
	node host.Host
}

func (d discoveryNotifee) HandlePeerFound(info peer.AddrInfo) {
	d.node.Connect(context.Background(), info)
}

func secondPeer(t *testing.T, discoveryName string) host.Host {
	host, _ := libp2p.New()
	discovery := mdns.NewMdnsService(host, discoveryName, &discoveryNotifee{node: host})
	_ = discovery.Start()
	t.Cleanup(func() {
		host.Close()
		discovery.Close()
	})
	return host
}

func TestNewSenderService(t *testing.T) {
	_, _, _, node, conf := setupEnvironment(t)

	sender := Sender{Node: node, Frequency: conf.Frequency}

	assert.NotNil(t, sender)
}

func TestSenderGetBandWidth(t *testing.T) {
	topic, subscr, ctx, node, conf := setupEnvironment(t)
	sender := Sender{Node: node, Frequency: conf.Frequency}
	discovery := mdns.NewMdnsService(node.Host, "discoveryTest", &discoveryNotifee{node: node.Host})
	_ = discovery.Start()
	peer2 := secondPeer(t, "discoveryTest")

	t.Cleanup(func() {
		node.Host.Close()
		discovery.Close()
		peer2.Close()
	})
	peer2info := metrics.NewSystemInfo("", peer2.ID(), "", "")

	sender.getBandwidth(peer2info, topic, ctx)
	message, _ := subscr.Next(ctx)

	assert.Contains(t, message.String(), "total_in")
	assert.Contains(t, message.String(), "total_out")
	assert.Contains(t, message.String(), "rate_in")
	assert.Contains(t, message.String(), "rate_in")
}
