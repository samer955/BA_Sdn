package subscriber

import (
	"context"
	"github.com/libp2p/go-libp2p-core/host"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"log"
	"time"
)

type PubSubService struct {
	psub *pubsub.PubSub
}

// NewPubSubService return a new PubSub Service using the GossipSub Service
func NewPubSubService(ctx context.Context, host host.Host) *PubSubService {

	ps, err := pubsub.NewGossipSub(ctx, host)
	if err != nil {
		time.Sleep(60 * time.Second)
		NewPubSubService(ctx, host)
	}
	return &PubSubService{psub: ps}
}

// JoinTopic allow the Peers to join a Topic on Pubsub
func (service *PubSubService) JoinTopic(room string) *pubsub.Topic {

	topic, err := service.psub.Join(room)
	if err != nil {
		log.Println("Error while subscribing in the Time-Topic")
	} else {
		log.Println("Subscribed on", room)
		log.Println("topicID", topic.String())
	}
	return topic
}

// Subscribe returns a new Subscription for the topic.
func (service *PubSubService) Subscribe(topic *pubsub.Topic) *pubsub.Subscription {

	subscribe, err := topic.Subscribe()

	if (err) != nil {
		log.Println("cannot subscribe to: ", topic.String())
	} else {
		log.Println("Subscribed to, " + subscribe.Topic())
	}
	return subscribe
}
