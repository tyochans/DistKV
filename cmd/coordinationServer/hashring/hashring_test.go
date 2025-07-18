package hashring

import (
	"testing"
)

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
	ring := NewHashRing(2, dummyHash)

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
