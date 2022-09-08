package node

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNodeCreateLibp2pNode(t *testing.T) {
	var node Node

	node.StartNode()

	assert.NotNil(t, node.Host)
}
