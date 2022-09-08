package node

import (
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/metrics"
	"log"
	"net"
	"sender-agent/config"
)

type Node struct {
	Host        host.Host
	Ip          string
	Role        string
	Network     string
	Bandcounter *metrics.BandwidthCounter
}

var conf = config.GetConfig()

func (n *Node) createLibp2pNode() {

	// create a new ibp2p Host that listens on a TCP port
	node, err := libp2p.New(libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/0"), libp2p.BandwidthReporter(n.Bandcounter))
	if err != nil {
		log.Fatal(err)
	}
	n.Host = node

}

func (n *Node) createBandCounter() {
	n.Bandcounter = metrics.NewBandwidthCounter()
}

// LocalIP get the node machine local IP address, based on the https://stackoverflow.com/questions/23558425/how-do-i-get-the-local-ip-address-in-go
func (n *Node) localIP() {

	ifaces, err := net.Interfaces()
	if err != nil {
		return
	}
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			return
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip.IsPrivate() {
				n.Ip = ip.String()
				return
			}
		}
	}
	n.Ip = ""

}

func (n *Node) getConfig() {

	n.Role = conf.Role
	n.Network = conf.Network

}

func (n *Node) StartNode() {

	n.localIP()
	n.createBandCounter()
	n.getConfig()
	n.createLibp2pNode()
	log.Println("created new node:", n.Host.ID().Pretty())

}
