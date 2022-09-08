package discovery

import (
	"context"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
	"log"
	"sender-agent/node"
	"sender-agent/service"
)

type discoveryNotifee struct {
	node host.Host
}

var PeerList []peer.AddrInfo
var pingTopic *pubsub.Topic

//used to publish the ping results on this topic after discovering a new peer
func SetPingTopic(topic *pubsub.Topic) {
	pingTopic = topic
}

//The node will be notificated when a new Peer is discovered
func (d *discoveryNotifee) HandlePeerFound(info peer.AddrInfo) {
	if d.node.ID().Pretty() != info.ID.Pretty() {
		log.Printf("discovered a new peer %s\n", info.ID.Pretty())

		err := d.node.Connect(context.Background(), info)
		if err != nil {
			log.Printf("unable to connect to Peer %s ", info.ID.Pretty())
			return
		}
		PeerList = append(PeerList, info)
		log.Printf("connected to Peer %s ", info.ID.Pretty())
		//once discovered a new Peer the local Host start to ping it. The result will be published in the pingTopic
		service.SendPing(context.Background(), d.node, info, pingTopic)
	}
}

func SetupDiscovery(node node.Node, discoveryName string) error {
	discovery := mdns.NewMdnsService(node.Host, discoveryName, &discoveryNotifee{node: node.Host})
	start := discovery.Start()
	return start
}
