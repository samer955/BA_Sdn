package metrics

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"os/exec"
	"testing"
)

var netstatOutput = "\nActive Connections\n" +
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

func TestNumbertcpQueue(t *testing.T) {
	var expected = 7

	got, _ := numberOfTcpQueue(netstatOutput)

	if got != expected {
		t.Fatal("wrong tcp queue number")
	}
}

func TestTcpSegmentsWindows(t *testing.T) {

	var netstatOutput = "" +
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

	var expReceived, expSent = 8011276, 7388244
	var gotReceived, gotSent, _ = numberOfSegmentsWindows(netstatOutput)

	if expReceived != gotReceived {
		t.Fatal("unexpected segments received")
	}
	if expSent != gotSent {
		t.Fatal("unexpected segments sent")
	}
}

func TestTcpSegmentsNumberLinux(t *testing.T) {

	var netstatOutput = "" +
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
		"\n    0 bad segments received"

	var exp_received, exp_sent = 5461, 7043
	var got_received, got_sent, _ = numbersOfSegmentsLinux(netstatOutput)

	if exp_received != got_received {
		t.Fatal("unexpected segments received")
	}
	if exp_sent != got_sent {
		t.Fatal("unexpected segments sent")
	}
}

//token from https://npf.io/2015/06/testing-exec-command/
func execCommandQueue(command string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestHelperProcess", "--", command}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	return cmd
}

//token from https://npf.io/2015/06/testing-exec-command/
func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	fmt.Fprintf(os.Stdout, netstatOutput)
	os.Exit(0)
}

func execCommandQueueError(command string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestHelperProcessError", "--", command}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	return cmd
}

func TestHelperProcessError(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	os.Exit(1)
}

func TestTcpQueueSize(t *testing.T) {
	tcp := NewTCPstatus("1.1.1.1")
	//basically mocking the exec.command of golang
	execCommand = execCommandQueue
	defer func() { execCommand = exec.Command }()

	tcp.TcpQueueSize()

	assert.Equal(t, tcp.QueueSize, 7)
}

func TestTcpQueueSizeWithError(t *testing.T) {
	tcp := NewTCPstatus("1.1.1.1")
	execCommand = execCommandQueueError
	defer func() { execCommand = exec.Command }()

	tcp.TcpQueueSize()

	assert.Equal(t, tcp.QueueSize, 0)

}
