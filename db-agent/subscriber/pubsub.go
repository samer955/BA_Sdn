package subscriber

import (
	"context"
	"github.com/libp2p/go-libp2p-core/host"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"log"
	"time"
)

var psub *pubsub.PubSub

// NewPubSubService return a new PubSub Service using the GossipSub Service
func NewPubSubService(ctx context.Context, host host.Host) *pubsub.PubSub {

	ps, err := pubsub.NewGossipSub(ctx, host)
	if err != nil {
		time.Sleep(60 * time.Second)
		NewPubSubService(ctx, host)
	}
	psub = ps
	return ps
}

// JoinTopic allow the Peers to join a Topic on Pubsub
func JoinTopic(room string) *pubsub.Topic {

	topic, err := psub.Join(room)
	if err != nil {
		log.Println("Error while subscribing in the Time-Topic")
	} else {
		log.Println("Subscribed on", room)
		log.Println("topicID", topic.String())
	}
	return topic
}

// Subscribe returns a new Subscription for the topic.
func Subscribe(topic *pubsub.Topic) *pubsub.Subscription {

	subscribe, err := topic.Subscribe()

	if (err) != nil {
		log.Println("cannot subscribe to: ", topic.String())
	} else {
		log.Println("Subscribed to, " + subscribe.Topic())
	}
	return subscribe
}
