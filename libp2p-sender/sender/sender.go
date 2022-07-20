package sender

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/beevik/ntp"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"libp2p-sender/discovery"
	"libp2p-sender/variables"
	"log"
	"time"
)

func SendTimeMessage(topic *pubsub.Topic, context context.Context, peer *variables.PeerInfo) {

	for {
		if len(discovery.PeerList) == 0 {
			continue
		}

		fmt.Println("sending time")

		peer.Time = TimeFromServer()

		//JSON encoding of peerInfo struct in order to send the data as []byte.
		peerInfoJson, _ := json.Marshal(peer)

		//public the Json content in the topic
		content := topic.Publish(context, peerInfoJson)

		if content != nil {
			log.Println("Error publishing content ", content.Error())
		}
		//wait 10 seconds before send another timestamp
		time.Sleep(5 * time.Second)
	}
}

// SendCpuInformation function will send information about the CPU
func SendCpuInformation(topic *pubsub.Topic, context context.Context, cpu *variables.Cpu) {
	for {
		if len(discovery.PeerList) == 0 {
			continue
		}
		fmt.Println("sending cpu")

		//Update every 10s CPU Usages in %
		updateCpuPercentage(cpu)

		//JSON encoding of cpu in order to send the data as []byte.
		jsonCpu, _ := json.Marshal(cpu)

		//public the data on the topic
		content := topic.Publish(context, jsonCpu)

		if content != nil {
			log.Println("Error publishing content ", content.Error())
		}

		time.Sleep(5 * time.Second)
	}
}

// SendRamInformation function will send information about the RAM
func SendRamInformation(topic *pubsub.Topic, context context.Context, ram *variables.Ram) {

	for {
		if len(discovery.PeerList) == 0 {
			continue
		}

		fmt.Println("sending ram")

		//Update every 10s RAM usages in %
		updateRamPercentage(ram)
		//JSON encoding of ram in order to send the data as []byte.
		jsonRam, _ := json.Marshal(ram)

		//public the data on the topic
		content := topic.Publish(context, jsonRam)

		if content != nil {
			log.Println("Error publishing content ", content.Error())
		}

		time.Sleep(5 * time.Second)
	}
}

func updateRamPercentage(ram *variables.Ram) {
	vmStat, _ := mem.VirtualMemory()
	ram.Usage = int(vmStat.UsedPercent)
	ram.Time = TimeFromServer()
}

func updateCpuPercentage(c *variables.Cpu) {
	cpuUsage, _ := cpu.Percent(0, false)
	c.Usage = int(cpuUsage[0])
	c.Time = TimeFromServer()
}

//This Function get the actual time from a server
func TimeFromServer() time.Time {
	now, err := ntp.Time("time.apple.com")
	if err != nil {
		fmt.Println(err)
	}
	return now
}
