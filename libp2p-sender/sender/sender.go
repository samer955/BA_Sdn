package sender

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/beevik/ntp"
	host2 "github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/p2p/protocol/ping"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/process"
	"libp2p-sender/variables"
	"math"
	"time"
)

//SendTimeMessage will send periodically a timestamp in order to calculate the latency in ms
//between sender and receiver
func SendTimeMessage(topic *pubsub.Topic, context context.Context, peer *variables.PeerInfo, peers []peer.AddrInfo) {

	for {
		if len(peers) == 0 {
			continue
		}

		fmt.Println("sending time")

		//Set the time when the message is sent
		peer.Time = TimeFromServer()

		err := publish(peer, context, topic)
		if err != nil {
			fmt.Println("Error publishing content ", err.Error())
		}

		//wait 10 seconds before send another timestamp
		time.Sleep(5 * time.Second)
	}
}

// SendCpuInformation function will send periodically information about the CPU
func SendCpuInformation(topic *pubsub.Topic, context context.Context, cpu *variables.Cpu, peers []peer.AddrInfo) {
	for {
		if len(peers) == 0 {
			continue
		}
		fmt.Println("sending cpu")

		//Update every 10s CPU Usages in %
		updateCpuPercentage(cpu)

		if cpu.Usage >= 80 {
			cpu.Processes = getProcessesCPU()
		}

		//publish the cpu data
		err := publish(cpu, context, topic)
		if err != nil {
			fmt.Println("Error publishing content ", err.Error())
		}

		//set the processes to null after publishing the data
		cpu.Processes = nil

		time.Sleep(15 * time.Second)
	}
}

// SendRamInformation function will send periodically information about the RAM
func SendRamInformation(topic *pubsub.Topic, context context.Context, ram *variables.Ram, peers []peer.AddrInfo) {

	for {
		if len(peers) == 0 {
			continue
		}

		fmt.Println("sending ram")

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

func SendPing(ctx context.Context, node host2.Host, peer peer.AddrInfo, topic *pubsub.Topic) {

	status := variables.PingStatus{
		Source_node: node.ID().Pretty(),
		Target_node: peer.ID.Pretty(),
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
			fmt.Println("pinged", peer.Addrs[0], "without success")
		}
		status.Time = TimeFromServer()

		//sendToDatabase(status)

		//Next Ping in 1 Min
		time.Sleep(10 * time.Second)
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

//TimeFromServer get the actual time from a remote server using the ntp Protocol
//The purpose is to synchronize the time between the VMs to avoid problems
func TimeFromServer() time.Time {
	now, err := ntp.Time("time.apple.com")
	if err != nil {
		fmt.Println(err)
	}
	return now
}

func getProcessesCPU() []variables.Process {

	var processList []variables.Process

	//get the actual processes running
	processes, err := process.Processes()

	if err != nil {
		fmt.Println("Unable to read processes")
		return nil
	}

	for _, proc := range processes {
		validateProcess(proc, &processList)
	}
	return processList
}

//validate a process: proof if the name is visible and the cpu % usage of it is more than 4%
func validateProcess(process *process.Process, list *[]variables.Process) {

	name, err := process.Name()
	if err == nil {
		perc, err := process.CPUPercent()
		if err == nil {
			if perc >= 4.0 {
				pro := variables.Process{
					Name:        name,
					Cpu_percent: math.Round(perc*100) / 100}

				*list = append(*list, pro)
			}
		}
	}
}

func sendToDatabase(status variables.PingStatus, context context.Context, topic *pubsub.Topic) {

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
