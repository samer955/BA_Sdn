package components

import (
	"github.com/stretchr/testify/assert"
	"testing"
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

	var netstat = "\nIPv4 Statistics\n" +
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
		"\n  Datagrams Failing Fragmentation    = 0\n  Fragments Created                  = 38488\n\nIPv6 Statistics\n\n  Packets Received                   = 0\n  Received Header Errors             = 0\n  Received Address Errors            = 0\n  Datagrams Forwarded                = 0\n  Unknown Protocols Received         = 0\n  Received Packets Discarded         = 0\n  Received Packets Delivered         = 336740\n  Output Requests                    = 339616\n  Routing Discards                   = 0\n  Discarded Output Packets           = 0\n  Output Packet No Route             = 0\n  Reassembly Required                = 0\n  Reassembly Successful              = 0\n  Reassembly Failures                = 0\n  Datagrams Successfully Fragmented  = 0\n  Datagrams Failing Fragmentation    = 0\n  Fragments Created                  = 0\n\nICMPv4 Statistics\n\n                            Received    Sent\n  Messages                  5873        9527\n  Errors                    0           0\n  Destination Unreachable   3979        7199\n  Time Exceeded             608         1\n  Parameter Problems        0           0\n  Source Quenches           0           0\n  Redirects                 0           0\n  Echo Replies              104         3\n  Echos                     148         2324\n  Timestamps                0           0\n  Timestamp Replies         0           0\n  Address Masks             0           0\n  Address Mask Replies      0           0\n  Router Solicitations      0           0\n  Router Advertisements     1034        0\n\nICMPv6 Statistics\n\n                            Received    Sent\n  Messages                  0           0\n  Errors                    0           0\n  Destination Unreachable   0           0\n  Packet Too Big            0           0\n  Time Exceeded             0           0\n  Parameter Problems        0           0\n  Echos                     0           0\n  Echo Replies              0           0\n  MLD Queries               0           0\n  MLD Reports               0           0\n  MLD Dones                 0           0\n  Router Solicitations      0           0\n  Router Advertisements     0           0\n  Neighbor Solicitations    0           0         \n  Neighbor Advertisements   0           0\n  Redirects                 0           0\n  Router Renumberings       0           0\n\nTCP Statistics for IPv4\n\n  Active Opens                        = 347883\n  Passive Opens                       = 86083\n  Failed Connection Attempts          = 13218\n  Reset Connections                   = 272367\n  Current Connections                 = 44\n  Segments Received                   = 8011276\n  Segments Sent                       = 7388244\n  Segments Retransmitted              = 34858\n\nTCP Statistics for IPv6\n\n  Active Opens                        = 390\n  Passive Opens                       = 172\n  No Ports              = 48499\n  Receive Errors        = 4\n  Datagrams Sent        = 1901984\n\nUDP Statistics for IPv6\n\n  Datagrams Received    = 1312791\n  No Ports              = 0\n  Receive Errors        = 0\n  Datagrams Sent        = 335866\n"

	var exp_received, exp_sent = 8011276, 7388244

	var got_received, got_sent, _ = numberOfSegmentsWindows(netstat)

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
