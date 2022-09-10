package service

import (
	"context"
	"db-agent/node"
	"db-agent/subscriber"
	"db-agent/variables"
	"encoding/json"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/stretchr/testify/assert"
	"testing"
)

type MockRepo struct {
	peer *variables.SystemInfo
	ram  *variables.Ram
	cpu  *variables.Cpu
	ping *variables.PingStatus
	tcp  *variables.TCPstatus
	band *variables.Bandwidth
	node string
}

func (m *MockRepo) Migrate() {
}

func (m *MockRepo) SaveSystemInfo(peer *variables.SystemInfo) *variables.SystemInfo {
	m.peer = peer
	return peer
}
func (m *MockRepo) SaveRamInfo(ram *variables.Ram) *variables.Ram {
	m.ram = ram
	return ram
}
func (m *MockRepo) SaveCpuIfo(cpu *variables.Cpu) *variables.Cpu {
	m.cpu = cpu
	return cpu
}
func (m *MockRepo) SavePingStatus(ping *variables.PingStatus) *variables.PingStatus {
	m.ping = ping
	return ping
}
func (m *MockRepo) SaveTCPstatus(tcp *variables.TCPstatus) *variables.TCPstatus {
	m.tcp = tcp
	return tcp
}
func (m *MockRepo) SaveBandwidth(band *variables.Bandwidth) *variables.Bandwidth {
	m.band = band
	return band
}
func (m *MockRepo) GetIpFromNode(node string) string {
	m.node = node
	return "1.1.1.1"
}

func setupEnvironment(t *testing.T) (*pubsub.Topic, *pubsub.Subscription, context.Context, node.Node) {
	const roomTest = "test"
	var node node.Node
	node.StartNode()
	ctx := context.Background()
	psub := subscriber.NewPubSubService(ctx, node)
	testTopic := psub.JoinTopic(roomTest)
	subsc := psub.Subscribe(testTopic)

	t.Cleanup(func() {
		node.Host.Close()
		ctx.Done()
		testTopic.Close()
		subsc.Cancel()
	})
	return testTopic, subsc, ctx, node
}

func TestReceiver_ReadRam(t *testing.T) {
	topic, subscr, ctx, node := setupEnvironment(t)
	ram := variables.Ram{Usage: 50}
	msg, _ := json.Marshal(ram)
	topic.Publish(ctx, msg)
	var mockrepo = &MockRepo{}
	receiver := Receiver{Node: node, Repository: mockrepo}

	receiver.readRamInformation(subscr, ctx)

	assert.Equal(t, mockrepo.ram.Usage, 50)
}

func TestReceiver_ReadCPU(t *testing.T) {
	topic, subscr, ctx, node := setupEnvironment(t)
	cpu := variables.Cpu{Model: "TEST_MODEL", Usage: 90}
	msg, _ := json.Marshal(cpu)
	topic.Publish(ctx, msg)
	var mockrepo = &MockRepo{}
	receiver := Receiver{Node: node, Repository: mockrepo}

	receiver.readCpuInformation(subscr, ctx)

	assert.Equal(t, mockrepo.cpu.Usage, 90)
	assert.Equal(t, mockrepo.cpu.Model, "TEST_MODEL")
}

func TestReceiver_ReadPingStatus(t *testing.T) {
	topic, subscr, ctx, node := setupEnvironment(t)
	pingStatus := variables.PingStatus{}
	msg, _ := json.Marshal(pingStatus)
	topic.Publish(ctx, msg)
	var mockrepo = &MockRepo{}
	receiver := Receiver{Node: node, Repository: mockrepo}

	receiver.readPingStatus(subscr, ctx)

	assert.Equal(t, mockrepo.ping.SourceIp, "1.1.1.1")
	assert.Equal(t, mockrepo.ping.TargetIp, "1.1.1.1")
	assert.Equal(t, mockrepo.ping.RTT, int64(0))
}

func TestReceiver_ReadBandwidth(t *testing.T) {
	topic, subscr, ctx, node := setupEnvironment(t)
	band := variables.Bandwidth{TotalIn: 10, TotalOut: 10, RateIn: 20, RateOut: 20}
	msg, _ := json.Marshal(band)
	topic.Publish(ctx, msg)
	var mockrepo = &MockRepo{}
	receiver := Receiver{Node: node, Repository: mockrepo}

	receiver.readBandwidth(subscr, ctx)

	assert.Equal(t, mockrepo.band.RateIn, 20)
	assert.Equal(t, mockrepo.band.RateOut, 20)
	assert.Equal(t, mockrepo.band.TotalIn, int64(10))
	assert.Equal(t, mockrepo.band.TotalOut, int64(10))
}

func TestReceiver_ReadTCPstatus(t *testing.T) {
	topic, subscr, ctx, node := setupEnvironment(t)
	tcp := variables.TCPstatus{QueueSize: 15, Sent: 20, Received: 20}
	msg, _ := json.Marshal(tcp)
	topic.Publish(ctx, msg)
	var mockrepo = &MockRepo{}
	receiver := Receiver{Node: node, Repository: mockrepo}

	receiver.readTCPstatus(subscr, ctx)

	assert.Equal(t, mockrepo.tcp.QueueSize, 15)
	assert.Equal(t, mockrepo.tcp.Sent, 20)
	assert.Equal(t, mockrepo.tcp.Received, 20)
}

func TestReceiver_ReadSystemInfo(t *testing.T) {
	topic, subscr, ctx, node := setupEnvironment(t)
	systemInfo := variables.SystemInfo{Hostname: "TESTNAME"}
	msg, _ := json.Marshal(systemInfo)
	topic.Publish(ctx, msg)
	var mockrepo = &MockRepo{}
	receiver := Receiver{Node: node, Repository: mockrepo}

	receiver.readSystemInfo(subscr, ctx)

	assert.NotNil(t, mockrepo.peer)
}
