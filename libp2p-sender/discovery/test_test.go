package discovery

import (
	"context"
	"fmt"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/p2p/discovery/mocks"
	"testing"
	"time"
)

type clock interface {
	Now() time.Time
}

type tempo struct {
	time time.Time
}

func (tempo) Now() time.Time {
	return time.Now()
}

func TestSetupDiscovery2(t *testing.T) {
	tempo := new(tempo)

	server := mocks.NewDiscoveryServer(tempo)

	node, _ := libp2p.New(libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/0"))

	client := mocks.NewDiscoveryClient(node, server)

	server.FindPeers("test", 10)
	client.FindPeers(context.Background(), "test")

	fmt.Println()

}
