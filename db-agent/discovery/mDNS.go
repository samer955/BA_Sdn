package discovery

import (
	"context"
	"fmt"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
	"log"
	"time"
)

const discoveryName = "discoveryRoom"

var PeerList []peer.AddrInfo

type discoveryNotifee struct {
	node host.Host
}

func (d *discoveryNotifee) HandlePeerFound(info peer.AddrInfo) {
	fmt.Printf("discovered a new peer %s\n", info.ID.Pretty())

	if d.node.ID().Pretty() != info.ID.Pretty() {
		d.node.Connect(context.Background(), info)
		PeerList = append(PeerList, info)

		log.Printf("connected to Peer %s ", info.ID.Pretty())
	}
}

func SetupDiscovery(node host.Host) error {
	discovery := mdns.NewMdnsService(node, discoveryName, &discoveryNotifee{node: node})
	start := discovery.Start()

	//If any error is returned try again in 1min
	if start != nil {
		time.Sleep(60 * time.Second)
		SetupDiscovery(node)
	}
	return start
}
