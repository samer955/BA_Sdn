package discovery

import (
	"db-agent/node"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestSetupDiscovery(t *testing.T) {
	host, _ := libp2p.New()
	var node = node.Node{Host: host}
	t.Cleanup(func() {
		node.Host.Close()
	})

	err := SetupDiscovery(node, "discoveryDB_Test")

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
	host, _ := libp2p.New()
	discovery := mdns.NewMdnsService(host, "discoveryRoomTest", &discoveryNotifee{node: host})
	_ = discovery.Start()

	t.Cleanup(func() {
		host.Close()
		discovery.Close()
	})

	secondPeer := secondPeer(t, "discoveryRoomTest")
	limit := time.Now()

	//wait till the other Peer is found, limit 4 seconds
	for {
		if time.Now().After(limit.Add(4 * time.Second)) {
			break
		}
		if len(host.Peerstore().Peers()) == 1 {
			continue
		}
		break
	}

	assert.Contains(t, host.Peerstore().Peers(), secondPeer.ID())
}
