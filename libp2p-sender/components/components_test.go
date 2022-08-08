package components

import (
	"errors"
	"fmt"
	"github.com/libp2p/go-libp2p/p2p/protocol/ping"
	"github.com/stretchr/testify/assert"
	"os/exec"
	"testing"
	"time"
)

func TestCpu(t *testing.T) {

	cp := cpuModel()
	assert.NotEqual(t, cp, "")

}

func TestCpu_UpdateCpuPercentage(t *testing.T) {
	cpu := NewCpu("1.1.1.1", "testID")

	cpu.UpdateCpuPercentage()
	assert.NotEqual(t, cpu.Usage, 0)
}

func TestNumbertcpQueue(t *testing.T) {

	var netstat = "\nActive Connections\n" +
		"\n  Proto  Local Address          Foreign Address        State" +
		"\n  TCP    0.0.0.0:135            0.0.0.0:0              LISTENING" +
		"\n  TCP    0.0.0.0:445            0.0.0.0:0              LISTENING" +
		"\n  TCP    0.0.0.0:3000           0.0.0.0:0              LISTENING" +
		"\n  TCP    127.0.0.1:6942         0.0.0.0:0              LISTENING" +
		"\n  TCP    127.0.0.1:49670        0.0.0.0:0              LISTENING" +
		"\n  TCP    127.0.0.1:50164        127.0.0.1:50165        ESTABLISHED" +
		"\n  TCP    127.0.0.1:50165        127.0.0.1:50164        ESTABLISHED" +
		"\n  TCP    127.0.0.1:50166        127.0.0.1:50167        ESTABLISHED" +
		"\n  TCP    127.0.0.1:50167        127.0.0.1:50166        ESTABLISHED" +
		"\n  TCP    172.20.10.3:50926      23.203.90.118:443      ESTABLISHED" +
		"\n  TCP    172.20.10.3:50983      23.203.90.118:443      ESTABLISHED" +
		"\n  TCP    172.20.10.3:50999      23.203.90.118:443      ESTABLISHED" +
		"\n  TCP    192.168.56.1:139       0.0.0.0:0              LISTENING" +
		"\n  TCP    192.168.208.1:139      0.0.0.0:0              LISTENING" +
		"\n  TCP    [::]:135               [::]:0                 LISTENING" +
		"\n  TCP    [::]:445               [::]:0                 LISTENING" +
		"\n  UDP    0.0.0.0:53             *:*" +
		"\n  UDP    0.0.0.0:123            *:*" +
		"\n  UDP    0.0.0.0:5050           *:*"

	var expected = 7
	got, _ := numberOfTcpQueue(netstat)

	if got != expected {
		t.Fatal("wrong tcp queue number")
	}
}

func TestTcpSegmentsWindows(t *testing.T) {

	var netstat = "" +
		"\nIPv4 Statistics\n" +
		"\n  Packets Received                   = 9147172" +
		"\n  Received Header Errors             = 0" +
		"\n  Received Address Errors            = 33" +
		"\n  Datagrams Forwarded                = 0" +
		"\n  Unknown Protocols Received         = 0" +
		"\n  Received Packets Discarded         = 57251" +
		"\n  Received Packets Delivered         = 9513095" +
		"\n  Output Requests                    = 7343447" +
		"\n  Routing Discards                   = 0" +
		"\n  Discarded Output Packets           = 3729" +
		"\n  Output Packet No Route             = 1154" +
		"\n  Reassembly Required                = 46940" +
		"\n  Reassembly Successful              = 23469" +
		"\n  Reassembly Failures                = 0" +
		"\n  Datagrams Successfully Fragmented  = 19244" +
		"\n  Datagrams Failing Fragmentation    = 0" +
		"\n  Fragments Created                  = 38488" +
		"\n" +
		"\nIPv6 Statistics\n" +
		"\n  Packets Received                   = 0" +
		"\n  Received Header Errors             = 0" +
		"\n  Received Address Errors            = 0" +
		"\n  Datagrams Forwarded                = 0" +
		"\n  Unknown Protocols Received         = 0" +
		"\n  Received Packets Discarded         = 0" +
		"\n  Received Packets Delivered         = 336740" +
		"\n  Output Requests                    = 339616" +
		"\n  Routing Discards                   = 0" +
		"\n  Discarded Output Packets           = 0" +
		"\n  Output Packet No Route             = 0" +
		"\n  Reassembly Required                = 0" +
		"\n  Reassembly Successful              = 0" +
		"\n  Reassembly Failures                = 0" +
		"\n  Datagrams Successfully Fragmented  = 0" +
		"\n  Datagrams Failing Fragmentation    = 0" +
		"\n  Fragments Created                  = 0" +
		"\n" +
		"\nICMPv4 Statistics\n" +
		"\n                            Received    Sent" +
		"\n  Messages                  5873        9527" +
		"\n  Errors                    0           0" +
		"\n  Destination Unreachable   3979        7199" +
		"\n  Time Exceeded             608         1" +
		"\n  Parameter Problems        0           0" +
		"\n  Source Quenches           0           0" +
		"\n  Redirects                 0           0" +
		"\n  Echo Replies              104         3" +
		"\n  Echos                     148         2324" +
		"\n  Timestamps                0           0" +
		"\n  Timestamp Replies         0           0" +
		"\n  Address Masks             0           0" +
		"\n  Address Mask Replies      0           0" +
		"\n  Router Solicitations      0           0" +
		"\n  Router Advertisements     1034        0" +
		"\n\nICMPv6 Statistics" +
		"\n" +
		"\n                            Received    Sent" +
		"\n  Messages                  0           0" +
		"\n  Errors                    0           0" +
		"\n  Destination Unreachable   0           0" +
		"\n  Packet Too Big            0           0" +
		"\n  Time Exceeded             0           0" +
		"\n  Parameter Problems        0           0" +
		"\n  Echos                     0           0" +
		"\n  Echo Replies              0           0" +
		"\n  MLD Queries               0           0" +
		"\n  MLD Reports               0           0" +
		"\n  MLD Dones                 0           0" +
		"\n  Router Solicitations      0           0" +
		"\n  Router Advertisements     0           0" +
		"\n  Neighbor Solicitations    0           0" +
		"\n  Neighbor Advertisements   0           0" +
		"\n  Redirects                 0           0" +
		"\n  Router Renumberings       0           0" +
		"\n" +
		"\nTCP Statistics for IPv4" +
		"\n" +
		"\n  Active Opens                        = 347883" +
		"\n  Passive Opens                       = 86083" +
		"\n  Failed Connection Attempts          = 13218" +
		"\n  Reset Connections                   = 272367" +
		"\n  Current Connections                 = 44" +
		"\n  Segments Received                   = 8011276" +
		"\n  Segments Sent                       = 7388244" +
		"\n  Segments Retransmitted              = 34858" +
		"\n" +
		"\nTCP Statistics for IPv6" +
		"\n" +
		"\n  Active Opens                        = 390" +
		"\n  Passive Opens                       = 172" +
		"\n  No Ports              = 48499" +
		"\n  Receive Errors        = 4" +
		"\n  Datagrams Sent        = 1901984" +
		"\n" +
		"\nUDP Statistics for IPv6" +
		"\n" +
		"\n  Datagrams Received    = 1312791" +
		"\n  No Ports              = 0" +
		"\n  Receive Errors        = 0" +
		"\n  Datagrams Sent        = 335866\n"

	var exp_received, exp_sent = 8011276, 7388244

	var got_received, got_sent, _ = numberOfSegmentsWindows(netstat)

	if exp_received != got_received {
		t.Fatal("unexpected segments received")
	}

	if exp_sent != got_sent {
		t.Fatal("unexpected segments sent")
	}
}

func TestTcpSegmentsNumberLinux(t *testing.T) {

	var netstat = "" +
		"IcmpMsg:" +
		"\n    InType3: 43" +
		"\n    InType8: 5" +
		"\n    InType9: 13" +
		"\n    OutType0: 5" +
		"\n    OutType3: 43" +
		"\nTcp:" +
		"\n    207 active connection openings" +
		"\n    0 passive connection openings" +
		"\n    2 failed connection attempts" +
		"\n    83 connection resets received" +
		"\n    26 connections established" +
		"\n    5461 segments received" +
		"\n    7043 segments sent out" +
		"\n    62 segments retransmitted" +
		"\n    0 bad segments received" +
		"\n    11 resets sent" +
		"\nUdpLite:" +
		"\nTcpExt:" +
		"\n    21 TCP sockets finished time wait in fast timer" +
		"\n    48 delayed acks sent" +
		"\n    Quick ack mode was activated 39 times" +
		"\n    1519 packet headers predicted" +
		"\n    1024 acknowledgments not containing data payload received" +
		"\n    497 predicted acknowledgments" +
		"\n    TCPSackRecovery: 4" +
		"\n    Detected reordering 6 times using SACK" +
		"\n    1 congestion windows recovered without slow start after partial ack" +
		"\n    TCPLostRetransmit: 19" +
		"\n    4 fast retransmits" +
		"\n    TCPTimeouts: 50" +
		"\n    TCPLossProbes: 23" +
		"\n    TCPLossProbeRecovery: 2" +
		"\n    TCPBacklogCoalesce: 2" +
		"\n    TCPDSACKOldSent: 41" +
		"\n    TCPDSACKRecv: 20" +
		"\n    3 connections reset due to unexpected data" +
		"\n    1 connections aborted due to timeout" +
		"\n    TCPDSACKIgnoredNoUndo: 6" +
		"\n    TCPSackShiftFallback: 9" +
		"\n    TCPRcvCoalesce: 447" +
		"\n    TCPOFOQueue: 32" +
		"\n    TCPSpuriousRtxHostQueues: 1" +
		"\n    TCPAutoCorking: 73" +
		"\n    TCPSynRetrans: 17" +
		"\n    TCPOrigDataSent: 2658" +
		"\n    TCPKeepAlive: 515" +
		"\n    TCPDelivered: 2865" +
		"\nIpExt:" +
		"\n    InNoRoutes: 1" +
		"\n    InMcastPkts: 855" +
		"\n    OutMcastPkts: 614" +
		"\n    InBcastPkts: 8" +
		"\n    OutBcastPkts: 4" +
		"\n    InOctets: 3989314" +
		"\n    OutOctets: 1319948" +
		"\n    InMcastOctets: 142731" +
		"\n    OutMcastOctets: 104424" +
		"\n    InBcastOctets: 872" +
		"\n    OutBcastOctets: 310" +
		"\n    InNoECTPkts: 8587"

	var exp_received, exp_sent = 5461, 7043

	var got_received, got_sent, _ = numbersOfSegmentsLinux(netstat)

	if exp_received != got_received {
		t.Fatal("unexpected segments received")
	}

	if exp_sent != got_sent {
		t.Fatal("unexpected segments sent")
	}
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

	assert.NotEqual(t, ram.Usage, 0)
}

func TestNewPingStatus(t *testing.T) {
	status := NewPingStatus("node_A", "node_B")

	assert.NotNil(t, status)
	assert.Equal(t, status.Source, "node_A")
	assert.Equal(t, status.Target, "node_B")

}

func TestCheckPingStatusPositive(t *testing.T) {

	status := NewPingStatus("node_A", "node_B")
	var pingDeadline = 10

	fmt.Println(status)
	result := ping.Result{Error: nil, RTT: 5 * time.Millisecond}

	status.SetPingStatus(result, &pingDeadline)

	assert.Equal(t, status.Alive, true)
	assert.Equal(t, status.RTT, int64(5))
	assert.NotEqual(t, status.UUID, "")
	assert.NotEqual(t, status.Time, (time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC)))

}

func TestCheckPingStatusNegative(t *testing.T) {

	status := NewPingStatus("node_A", "node_B")
	var pingDeadline = 10

	fmt.Println(status)
	result := ping.Result{Error: errors.New("any Error"), RTT: 5 * time.Millisecond}

	status.SetPingStatus(result, &pingDeadline)

	assert.Equal(t, status.Alive, false)
	assert.Equal(t, status.RTT, int64(0))
	assert.NotEqual(t, status.UUID, "")
	assert.NotEqual(t, status.Time, (time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC)))
	assert.Equal(t, pingDeadline, 9)

}

func TestGetNumberOfOnlineUseLinux(t *testing.T) {

	var outputWithZeroUser = ""

	var outputWithUsers = "" +
		"ubuntu   pts/0        2022-08-08 08:12 (62.84.220.226)\n" +
		"ubuntu   pts/1        2022-08-08 09:56 (62.84.220.226)"

	zeroUser := outputToIntLinux(outputWithZeroUser)
	onlineUser := outputToIntLinux(outputWithUsers)

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

	zeroUser := outputToIntWindows(outputWithZeroUser)
	twoUsers := outputToIntWindows(outputWithTwoUsers)

	assert.Equal(t, zeroUser, 0)
	assert.Equal(t, twoUsers, 2)
}

func Test(t *testing.T) {
	out, err := exec.Command("query", "user").Output()

	if err != nil {
		fmt.Println("unableToRead")
	}

	fmt.Println(string(out))
}
