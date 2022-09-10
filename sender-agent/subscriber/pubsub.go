package subscriber

import (
	"context"
	"encoding/json"
	"github.com/libp2p/go-libp2p-pubsub"
	"log"
	"sender-agent/node"
)

type PubSubService struct {
	psub *pubsub.PubSub
}

// NewPubSubService return a new PubSub Service using the GossipSub Service
func NewPubSubService(ctx context.Context, node node.Node) *PubSubService {
	ps, err := pubsub.NewGossipSub(ctx, node.Host)
	if err != nil {
		log.Println("unable to create the pubsub service")
	}
	return &PubSubService{psub: ps}
}

// JoinTopic allow the Peers to join a Topic on Pubsub
func (service *PubSubService) JoinTopic(room string) *pubsub.Topic {

	topic, err := service.psub.Join(room)
	if err != nil {
		log.Println("Error while subscribing in", room)
	} else {
		log.Println("Joined room:", room)
	}

	return topic

}

// Subscribe returns a new Subscription for the topic.
func (service *PubSubService) Subscribe(topic *pubsub.Topic) *pubsub.Subscription {

	subscribe, err := topic.Subscribe()

	if err != nil {
		log.Println("cannot subscribe to: ", topic.String())
	} else {
		log.Println("Subscribed to topic: " + subscribe.Topic())
	}
	return subscribe
}

func Publish(object interface{}, context context.Context, topic *pubsub.Topic) error {

	//JSON encoding of cpu in order to send the data as []byte.
	msgBytes, err := json.Marshal(object)

	if err != nil {
		log.Println("cannot convert to Bytes ", object)
	}
	//public the data in the topic
	return topic.Publish(context, msgBytes)
}
