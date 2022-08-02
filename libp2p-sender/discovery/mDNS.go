package discovery

import (
	"context"
	"fmt"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
	"libp2p-sender/service"
	"log"
	"time"
)

type discoveryNotifee struct {
	node host.Host
}

var PeerList []peer.AddrInfo
var pingTopic *pubsub.Topic

func SetPingTopic(topic *pubsub.Topic) {
	pingTopic = topic
}

//The node will be notificated when a new Peer is discovered
func (d *discoveryNotifee) HandlePeerFound(info peer.AddrInfo) {
	fmt.Printf("discovered a new peer %s\n", info.ID.Pretty())

	if d.node.ID().Pretty() != info.ID.Pretty() {
		d.node.Connect(context.Background(), info)
		PeerList = append(PeerList, info)

		fmt.Println(d.node.Peerstore().Peers())

		log.Printf("connected to Peer %s ", info.ID.Pretty())

		service.SendPing(context.Background(), d.node, info, pingTopic)
	}
}

func SetupDiscovery(node host.Host, discoveryName string) error {
	discovery := mdns.NewMdnsService(node, discoveryName, &discoveryNotifee{node: node})

	start := discovery.Start()

	//If any error is returned try again in 1min
	if start != nil {
		time.Sleep(60 * time.Second)
		SetupDiscovery(node, discoveryName)
	}
	return start
}
