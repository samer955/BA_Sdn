package discovery

import (
	"context"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
	"github.com/stretchr/testify/assert"
	"sender-agent/node"
	"sender-agent/subscriber"
	"testing"
	"time"
)

func setupEnvironment(t *testing.T) (node.Node, *pubsub.Topic) {
	const roomTest = "test"
	ctx := context.Background()
	host, _ := libp2p.New()
	testNode := node.Node{Host: host}
	psub := subscriber.NewPubSubService(ctx, testNode)
	topic := psub.JoinTopic(roomTest)
	SetPingTopic(topic)

	t.Cleanup(func() {
		testNode.Host.Close()
		ctx.Done()
	})
	return testNode, topic
}

func TestSetPingTopic(t *testing.T) {
	_, _ = setupEnvironment(t)
	assert.NotNil(t, pingTopic)
}

func TestSetupDiscovery(t *testing.T) {
	testNode, _ := setupEnvironment(t)

	err := SetupDiscovery(testNode, "test_0")

	assert.Nil(t, err)
}

func secondPeer(t *testing.T, discoveryName string) host.Host {
	node, _ := libp2p.New()
	discovery := mdns.NewMdnsService(node, discoveryName, &discoveryNotifee{node: node})
	_ = discovery.Start()
	t.Cleanup(func() {
		node.Close()
		discovery.Close()
	})
	return node
}

func TestDiscoveryHandlePeerFound(t *testing.T) {
	testNode, _ := setupEnvironment(t)
	SetupDiscovery(testNode, "discoveryRoomTest")
	secondPeer := secondPeer(t, "discoveryRoomTest")
	limit := time.Now()

	//wait till the other Peer is found, limit 4 seconds
	for {
		if time.Now().After(limit.Add(4 * time.Second)) {
			break
		}
		if len(testNode.Host.Peerstore().Peers()) == 1 {
			continue
		}
		break
	}

	assert.Contains(t, testNode.Host.Peerstore().Peers(), secondPeer.ID())
}
