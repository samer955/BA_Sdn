package service

import (
	"context"
	"encoding/json"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"libp2p-sender/metrics"
	"libp2p-sender/subscriber"
	"testing"
	"time"
)

type Mocking struct {
	mock.Mock
}

func setup(node host.Host, roomtest string) (*pubsub.Topic, *pubsub.Subscription, context.Context) {

	ctx := context.Background()
	_ = subscriber.NewPubSubService(ctx, node)
	testTopic := subscriber.JoinTopic(roomtest)
	subsc := subscriber.Subscribe(testTopic)

	return testTopic, subsc, ctx

}

func TestSendPeerInfo(t *testing.T) {

	node, _ := libp2p.New(libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/0"))
	peer_info := metrics.NewPeerInfo("1.1.1.1", "test_ID", "sender")
	topic, subscr, ctx := setup(node, "test")

	t.Cleanup(func() {
		node.Close()
		ctx.Done()
		topic.Close()
		subscr.Cancel()
	})

	sendPeerInfo(topic, ctx, peer_info, nil)
	message, _ := subscr.Next(ctx)

	peerToBytesConverted, _ := json.Marshal(peer_info)

	assert.Equal(t, message.Data, peerToBytesConverted)
	assert.NotEqual(t, peer_info.UUID, "")
	assert.NotEqual(t, peer_info.Time, time.Time{})
}

func TestSendCpuInfo(t *testing.T) {

	node, _ := libp2p.New(libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/0"))
	peer_cpu := metrics.NewCpu("1.1.1.1", "test_ID")
	topic, subscr, ctx := setup(node, "test")

	t.Cleanup(func() {
		node.Close()
		ctx.Done()
		topic.Close()
		subscr.Cancel()
	})

	sendCpuInfo(topic, nil, ctx, peer_cpu)
	message, _ := subscr.Next(ctx)

	peerToBytesConverted, _ := json.Marshal(peer_cpu)

	assert.Equal(t, message.Data, peerToBytesConverted)
	assert.NotEqual(t, peer_cpu.UUID, "")
	assert.NotEqual(t, peer_cpu.Usage, 0)
}

func TestSendRamInfo(t *testing.T) {

	node, _ := libp2p.New(libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/0"))
	peer_ram := metrics.NewRam("1.1.1.1", "test_ID")
	topic, subscr, ctx := setup(node, "test")

	t.Cleanup(func() {
		node.Close()
		ctx.Done()
		topic.Close()
		subscr.Cancel()
	})

	sendRamInfo(topic, nil, ctx, peer_ram)
	message, _ := subscr.Next(ctx)

	peerToBytesConverted, _ := json.Marshal(peer_ram)

	assert.Equal(t, message.Data, peerToBytesConverted)
	assert.NotEqual(t, peer_ram.UUID, "")
	assert.NotEqual(t, peer_ram.Usage, 0)
}
