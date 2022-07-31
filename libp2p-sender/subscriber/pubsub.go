package subscriber

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-pubsub"
	"log"
)

var psub *pubsub.PubSub

// NewPubSubService return a new PubSub Service using the GossipSub Service
func NewPubSubService(ctx context.Context, host host.Host) *pubsub.PubSub {

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
	return topic
}

func Subscribe(topic *pubsub.Topic) *pubsub.Subscription {

	subscribe, err := topic.Subscribe()

	if err != nil {
		log.Println("cannot subscribe to: ", topic.String())
	} else {
		log.Println("Subscribed to, " + subscribe.Topic())
	}
	return subscribe
}

func Publish(object interface{}, context context.Context, topic *pubsub.Topic) error {

	//JSON encoding of cpu in order to send the data as []byte.
	msgBytes, err := json.Marshal(object)

	if err != nil {
		fmt.Println("cannot convert to Bytes ", object)
	}
	//public the data in the topic
	return topic.Publish(context, msgBytes)
}
