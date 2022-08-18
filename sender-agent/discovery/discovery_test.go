package discovery

import (
	"context"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
	"github.com/stretchr/testify/assert"
	"libp2p-sender/subscriber"
	"testing"
)

func TestSetPingTopic(t *testing.T) {

	const roomTest = "test"
	host, _ := libp2p.New(libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/0"))
	ctx := context.Background()

	t.Cleanup(func() {
		host.Close()
		ctx.Done()
	})

	_ = subscriber.NewPubSubService(ctx, host)
	topic := subscriber.JoinTopic(roomTest)

	SetPingTopic(topic)

	assert.NotNil(t, pingTopic)
}

func TestSetupDiscovery(t *testing.T) {

	node, _ := libp2p.New(libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/0"))
	err := SetupDiscovery(node, "test_0")

	t.Cleanup(func() {
		node.Close()
	})

	assert.Nil(t, err)
}

func secondPeer(t *testing.T, discoveryName string) host.Host {
	node, _ := libp2p.New(libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/0"))
	discovery := mdns.NewMdnsService(node, discoveryName, &discoveryNotifee{node: node})
	_ = discovery.Start()
	t.Cleanup(func() {
		node.Close()
		discovery.Close()
	})
	return node
}

func TestDiscoveryNotifee_HandlePeerFound(t *testing.T) {

	host, _ := libp2p.New(libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/0"))
	SetupDiscovery(host, "discoveryRoomTest")

	t.Cleanup(func() {
		host.Close()
	})

	node := secondPeer(t, "discoveryRoomTest")

	//wait till the other Peer is found
	for {
		if len(PeerList) == 0 {
			continue
		}
		break
	}

	assert.Contains(t, host.Peerstore().Peers(), node.ID())
}
