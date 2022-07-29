package sender

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/beevik/ntp"
	"github.com/google/uuid"
	host2 "github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/p2p/protocol/ping"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"libp2p-sender/variables"
	"time"
)

//SendPeerInfo will send system Information of this Peer and periodically a timestamp
//in order to calculate the latency in ms between sender and receiver
func SendPeerInfo(topic *pubsub.Topic, context context.Context, peer *variables.PeerInfo, list *[]peer.AddrInfo) {

	for {
		if len(*list) == 0 {
			continue
		}

		fmt.Println("sending time")

		//Set the time when the message is sent
		peer.Time = TimeFromServer()
		peer.UUID = uuid.New().String()

		err := publish(peer, context, topic)
		if err != nil {
			fmt.Println("Error publishing content ", err.Error())
		}

		//wait 20 seconds before send another timestamp
		time.Sleep(20 * time.Second)
	}
}

// SendCpuInformation function will send periodically information about the CPU
func SendCpuInformation(topic *pubsub.Topic, context context.Context, cpu *variables.Cpu, peers *[]peer.AddrInfo) {
	for {
		if len(*peers) == 0 {
			continue
		}
		fmt.Println("sending cpu")

		cpu.UUID = uuid.New().String()
		//Update every 10s CPU Usages in %
		updateCpuPercentage(cpu)

		//publish the cpu data
		err := publish(cpu, context, topic)

		if err != nil {
			fmt.Println("Error publishing content ", err.Error())
		}

		time.Sleep(15 * time.Second)
	}
}

// SendRamInformation function will send periodically information about the actual RAM Percentage
func SendRamInformation(topic *pubsub.Topic, context context.Context, ram *variables.Ram, peers *[]peer.AddrInfo) {
	for {
		if len(*peers) == 0 {
			continue
		}
		fmt.Println("sending ram")

		ram.UUID = uuid.New().String()
		//Update every 10s RAM usages in %
		updateRamPercentage(ram)

		//publish the ram data
		err := publish(ram, context, topic)
		if err != nil {
			fmt.Println("Error publishing content ", err.Error())
		}
		time.Sleep(5 * time.Second)
	}
}

//SendPing function send a Ping every 60s to the discovered Peer
//This function is used in order to reach the others nodes in an active way getting
//a bool if the ping was successfully (true if a node is reachable, false if not) and an RTT in ms
func SendPing(ctx context.Context, node host2.Host, peer peer.AddrInfo, topic *pubsub.Topic) {

	status := variables.PingStatus{
		Source: node.ID().Pretty(),
		Target: peer.ID.Pretty(),
	}
	//The Ping function return a channel that still open till the context is alive
	ch := ping.Ping(ctx, node, peer.ID)

	for {
		res := <-ch

		if res.Error == nil {
			status.Alive = true
			status.RTT = res.RTT.Milliseconds()
			fmt.Println("pinged", peer.Addrs[0], "in", res.RTT)
		} else {
			status.Alive = false
			status.RTT = 0
			fmt.Println("pinged", peer.Addrs[0], "without success", res.Error)
		}
		status.Time = TimeFromServer()
		status.UUID = uuid.New().String()

		//publish the status of the Ping in the topic
		publish(status, ctx, topic)

		//Next Ping in 1 Min
		time.Sleep(59 * time.Second)
	}
}

//Get the actual RAM Percentage from the system
func updateRamPercentage(ram *variables.Ram) {
	ram.Time = TimeFromServer()
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		fmt.Println("Unable to get Memory Info")
		ram.Usage = 00
		return
	}
	ram.Usage = int(vmStat.UsedPercent)
}

//Get the actual CPU Percentage from the system
func updateCpuPercentage(c *variables.Cpu) {
	c.Time = TimeFromServer()
	cpuUsage, err := cpu.Percent(0, false)
	if err != nil {
		fmt.Println("Unable to get Cpu Percentage")
		c.Usage = 00
		return
	}
	c.Usage = int(cpuUsage[0])

}

//TimeFromServer get the actual time from a remote server using the ntp Protocol
//The purpose is to synchronize the time between the VMs to avoid problems
func TimeFromServer() time.Time {
	now, err := ntp.Time("time.apple.com")
	if err != nil {
		fmt.Println(err)
	}
	return now
}

func publish(object interface{}, context context.Context, topic *pubsub.Topic) error {

	//JSON encoding of cpu in order to send the data as []byte.
	msgBytes, err := json.Marshal(object)

	if err != nil {
		fmt.Println("cannot convert to Bytes ", object)
	}
	//public the data in the topic
	return topic.Publish(context, msgBytes)
}
