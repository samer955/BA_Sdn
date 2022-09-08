package discovery

import (
	"context"
	"db-agent/node"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
	"log"
)

type discoveryNotifee struct {
	node host.Host
}

var PeerList []peer.AddrInfo

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
	}

}

func SetupDiscovery(node node.Node, discoveryName string) error {

	discovery := mdns.NewMdnsService(node.Host, discoveryName, &discoveryNotifee{node: node.Host})
	start := discovery.Start()
	return start

}
