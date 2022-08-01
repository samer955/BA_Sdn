package subscriber

import (
	"context"
	"encoding/json"
	"github.com/libp2p/go-libp2p"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewPubSubService(t *testing.T) {

	ctx := context.Background()
	node, _ := libp2p.New(libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/0"))

	ps := NewPubSubService(ctx, node)

	assert.NotEqual(t, ps, nil)
}

func TestJoinTopic(t *testing.T) {

	const roomtest = "test"

	ctx := context.Background()
	node, _ := libp2p.New(libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/0"))
	_ = NewPubSubService(ctx, node)

	testTopic := JoinTopic(roomtest)

	assert.Equal(t, testTopic.String(), roomtest)
}

func TestSubscribe(t *testing.T) {

	const roomtest = "test"

	ctx := context.Background()
	node, _ := libp2p.New(libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/0"))
	ps := NewPubSubService(ctx, node)

	testTopic := JoinTopic(roomtest)
	_ = Subscribe(testTopic)

	assert.Contains(t, ps.GetTopics(), testTopic.String())
}

func TestPublish(t *testing.T) {

	type Message struct {
		Data string
	}
	const roomtest = "test"
	helloMessage := new(Message)
	helloMessage.Data = "Hello World"

	ctx := context.Background()
	node, _ := libp2p.New(libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/0"))
	_ = NewPubSubService(ctx, node)
	testTopic := JoinTopic(roomtest)
	subscr := Subscribe(testTopic)
	Publish(helloMessage, ctx, testTopic)

	//read the message published
	received, _ := subscr.Next(ctx)

	//unmarshal message data
	receivedMess := new(Message)
	json.Unmarshal(received.Data, receivedMess)

	assert.Equal(t, receivedMess.Data, helloMessage.Data)
}
