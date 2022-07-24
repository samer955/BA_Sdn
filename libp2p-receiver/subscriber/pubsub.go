package subscriber

import (
	"context"
	"github.com/libp2p/go-libp2p-core/host"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"log"
)

var psub *pubsub.PubSub

//map of topics
var topics map[string]*pubsub.Topic

// NewPubSubService return a new PubSub Service using the GossipSub Service
func NewPubSubService(ctx context.Context, host host.Host) *pubsub.PubSub {

	topics = make(map[string]*pubsub.Topic)

	ps, err := pubsub.NewGossipSub(ctx, host)
	if err != nil {
		panic(err)
	}
	psub = ps
	return ps
}

func JoinTopic(room string) *pubsub.Topic {

	topic, err := psub.Join(room)
	if err != nil {
		log.Println("Error while subscribing in the Time-Topic")
	} else {
		log.Println("Subscribed on", room)
		log.Println("topicID", topic.String())
	}
	//save topic on a map
	topics[room] = topic
	return topic
}

func Subscribe(topic *pubsub.Topic) *pubsub.Subscription {

	subscribe, err := topic.Subscribe()

	if (err) != nil {
		log.Println("cannot subscribe to: ", topic.String())
	} else {
		log.Println("Subscribed to, " + subscribe.Topic())
	}
	return subscribe
}
