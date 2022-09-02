package discovery

import (
	"context"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
	"log"
	"time"
)

var PeerList []peer.AddrInfo

type discoveryNotifee struct {
	node host.Host
}

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

func SetupDiscovery(node host.Host, discoveryName string) error {
	discovery := mdns.NewMdnsService(node, discoveryName, &discoveryNotifee{node: node})
	start := discovery.Start()

	//If any error is returned try again in 1min
	if start != nil {
		time.Sleep(60 * time.Second)
		SetupDiscovery(node, "")
	}
	return start
}
