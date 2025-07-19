package hashring

import (
	"hash/fnv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func simpleHash(data []byte) uint64 {
	h := fnv.New64a()
	h.Write(data)
	return h.Sum64()
}

// Dummy hash function for predictability
func dummyHash(data []byte) uint64 {
	switch string(data) {
	case "node-A#0": return 10
	case "node-A#1": return 20
	case "node-B#0": return 50
	case "node-B#1": return 70
	case "apple": return 15
	case "grape": return 60
	}
	return 0
}

func TestGetNode_Mapping(t *testing.T) {
	ring := NewHashRing(2,2,dummyHash)

	ring.AddNode("node-A")
	ring.AddNode("node-B")

	node1 := ring.GetNode("apple")
	if node1 != "node-A" {
		t.Errorf("Expected apple → node-A, got %s", node1)
	}

	node2 := ring.GetNode("grape")
	if node2 != "node-B" {
		t.Errorf("Expected grape → node-B, got %s", node2)
	}
}

func TestGetReplicas_Basic(t *testing.T) {
	hr := NewHashRing(3, 3, simpleHash)

	hr.AddNode("127.0.0.1:9001")
	hr.AddNode("127.0.0.1:9002")
	hr.AddNode("127.0.0.1:9003")

	replicas := hr.GetReplicas("user42")

	assert.Len(t, replicas, 3, "should return 3 distinct nodes")
	assert.ElementsMatch(t, []string{
		"127.0.0.1:9001", "127.0.0.1:9002", "127.0.0.1:9003",
	}, replicas)
}

func TestGetReplicas_NoDuplicateRealNodes(t *testing.T) {
	hr := NewHashRing(5,3, simpleHash) // 5 vnodes per real node

	hr.AddNode("127.0.0.1:9001")
	replicas := hr.GetReplicas("key1")

	assert.Len(t, replicas, 1, "should return only 1 real node due to only 1 registered")
	assert.Equal(t, "127.0.0.1:9001", replicas[0])
}

func TestGetReplicas_WrapAround(t *testing.T) {
	hr := NewHashRing(1, 3,simpleHash) // 1 vnode per real node to control placement

	// Intentionally register out of order
	hr.AddNode("9003")
	hr.AddNode("9001")
	hr.AddNode("9002")

	replicas := hr.GetReplicas("zzzzzzzzzz") // hashes high, past last node

	assert.Len(t, replicas, 3)
	assert.ElementsMatch(t, []string{"9001", "9002", "9003"}, replicas)
}
