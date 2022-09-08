package metrics

import (
	"errors"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p/p2p/protocol/ping"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCpu(t *testing.T) {
	cpu := cpuModel()

	assert.NotEqual(t, cpu, "")
}

func TestCpu_UpdateCpuPercentage(t *testing.T) {
	cpu := NewCpu("1.1.1.1", "testID")

	cpu.UpdateCpuPercentage()

	assert.NotEqual(t, cpu, nil)
	assert.NotEqual(t, cpu.Usage, nil)
}

func TestNewRam(t *testing.T) {
	ip := "1.1.1.1"
	node := "testNode"

	ram := NewRam(ip, node)

	assert.NotEqual(t, ram, nil)
	assert.Equal(t, ram.Ip, ip)
	assert.Equal(t, ram.Id, node)
}

func TestRam_UpdateRamPercentage(t *testing.T) {
	ip := "1.1.1.1"
	node := "testNode"
	ram := NewRam(ip, node)

	ram.UpdateRamPercentage()

	assert.NotEqual(t, ram.Usage, nil)
}

func TestNewPingStatus(t *testing.T) {
	status := NewPingStatus("node_A", "node_B")

	assert.NotNil(t, status)
	assert.Equal(t, status.Source, "node_A")
	assert.Equal(t, status.Target, "node_B")

}

func TestCheckPingStatusPositive(t *testing.T) {
	status := NewPingStatus("node_A", "node_B")
	var actualNegativePing = 0
	result := ping.Result{Error: nil, RTT: 5 * time.Millisecond}

	status.SetPingStatus(result, &actualNegativePing)

	assert.Equal(t, status.Alive, true)
	assert.Equal(t, status.RTT, int64(5))
	assert.Equal(t, actualNegativePing, 0)

}

func TestCheckPingStatusNegative(t *testing.T) {
	status := NewPingStatus("node_A", "node_B")
	var actualnegativePing = 0
	result := ping.Result{Error: errors.New("any Error"), RTT: 0 * time.Millisecond}

	status.SetPingStatus(result, &actualnegativePing)

	assert.Equal(t, status.Alive, false)
	assert.Equal(t, status.RTT, int64(0))
	assert.Equal(t, actualnegativePing, 1)

}

func TestGetNumberOfOnlineUseLinux(t *testing.T) {
	var outputWithZeroUser = ""
	var outputWithUsers = "" +
		"ubuntu   pts/0        2022-08-08 08:12 (62.84.220.226)\n" +
		"ubuntu   pts/1        2022-08-08 09:56 (62.84.220.226)"

	zeroUser := outputToIntUserLinux(outputWithZeroUser)
	onlineUser := outputToIntUserLinux(outputWithUsers)

	assert.Equal(t, zeroUser, 0)
	assert.Equal(t, onlineUser, 2)
}

func TestGetNumberOfOnlineUserWindows(t *testing.T) {
	var outputWithZeroUser = ""
	var outputWithTwoUsers = "" +
		" USERNAME              SESSIONNAME        ID  STATE   IDLE TIME  LOGON TIME\n" +
		">s.osman               console            31  Active      13:53  8/8/2022 10:06 AM\n" +
		">s.osman               console            32  Active      13:54  8/8/2022 10:16 AM\n" +
		">s.osman               console            33  Offline     13:55  8/8/2022 10:17 AM"

	zeroUser := outputToIntUserWin(outputWithZeroUser)
	twoUsers := outputToIntUserWin(outputWithTwoUsers)

	assert.Equal(t, zeroUser, 0)
	assert.Equal(t, twoUsers, 2)
}

func TestNewBandWidth(t *testing.T) {
	ip := "1.1.1.1"
	nodeId := "testNode"

	band := NewBandWidth(ip, nodeId)

	assert.NotEqual(t, band, nil)
	assert.Equal(t, band.Source, ip)
	assert.Equal(t, band.Id, nodeId)
}

func TestNewPeerInfo(t *testing.T) {
	ip := "1.1.1.1"
	peerID := peer.ID("test_node")
	role := "TEST_SENDER"
	network := "HOME"

	peer := NewSystemInfo(ip, peerID, role, network)

	assert.NotEqual(t, peer, nil)
	assert.NotEqual(t, peer.Hostname, "")
	assert.NotEqual(t, peer.OS, "")
	assert.NotEqual(t, peer.Version, "")
	assert.NotEqual(t, peer.Platform, "")
	assert.Equal(t, peer.Network, network)
	assert.Equal(t, peer.Role, role)
}
