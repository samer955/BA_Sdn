package node

import (
	"github.com/stretchr/testify/assert"
	"sender-agent/config"
	"strings"
	"testing"
)

func TestLocalIP(t *testing.T) {
	var node Node
	node.localIP()
	isPrivate := false
	ip := node.Ip

	if ip == "" {
		return
	}
	if strings.HasPrefix(ip, "10.") ||
		strings.HasPrefix(ip, "172.") ||
		strings.HasPrefix(ip, "192.") {
		isPrivate = true
	}
	assert.NotEqual(t, ip, "")
	assert.Equal(t, isPrivate, true)
}

func TestNodeCreateBandCounter(t *testing.T) {
	var node Node

	node.createBandCounter()

	assert.NotNil(t, node.Bandcounter)
}

func TestCreateLibp2pHost(t *testing.T) {
	var node Node

	node.createLibp2pNode()

	assert.NotNil(t, node.Host)
}

var mockConfig = config.Config{Network: "TEST_HOME", Role: "TEST_NODE"}

func GetMockConfig() config.Config {
	return mockConfig
}

func TestConfig(t *testing.T) {
	var node Node
	conf = GetMockConfig()
	defer func() { conf = config.GetConfig() }()

	node.getConfig()

	assert.Equal(t, node.Role, "TEST_NODE")
	assert.Equal(t, node.Network, "TEST_HOME")
}

func TestNode_GetNodeReady(t *testing.T) {
	var node Node
	conf = GetMockConfig()
	defer func() { conf = config.GetConfig() }()

	node.StartNode()

	assert.NotNil(t, node.Bandcounter)
	assert.NotNil(t, node.Host)
	assert.NotNil(t, node.Ip)
	assert.Equal(t, node.Network, "TEST_HOME")
	assert.Equal(t, node.Role, "TEST_NODE")
}
