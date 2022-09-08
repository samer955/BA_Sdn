package node

import (
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	"log"
)

type Node struct {
	Host host.Host
}

func (n *Node) StartNode() {

	// create a new ibp2p Host that listens on a TCP port
	node, err := libp2p.New(libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/0"))
	if err != nil {
		log.Println("unable to start the node")
		log.Fatal(err)
	}
	log.Println("created a new node:", node.ID().Pretty())
	n.Host = node

}
