package main

import (
	"DistKV/cmd/coordinationServer/hashring"
	"bufio"
	"net"
	"strings"
	"testing"
	"time"
)


func startMockWorker(t *testing.T, addr string, outputChan chan string) {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		t.Fatalf("Failed to start mock worker at %s: %v", addr, err)
	}

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				continue
			}
			go func(c net.Conn) {
				defer c.Close()
				reader := bufio.NewReader(c)
				line, _ := reader.ReadString('\n')
				outputChan <- strings.TrimSpace(line)
			}(conn)
		}
	}()
}

func TestReplicationForwarding(t *testing.T) {
	// Setup 3 mock worker servers
	workerPorts := []string{":9101", ":9102", ":9103"}
	outputChan := make(chan string, 3)

	for _, port := range workerPorts {
		startMockWorker(t, port, outputChan)
	}

	// Build hash ring
	ring := hashring.NewHashRing(3, 3, fnvHash)
	ring.AddNode("localhost:9101")
	ring.AddNode("localhost:9102")
	ring.AddNode("localhost:9103")

	// Trigger replication
	replicas := ring.GetReplicas("testKey")
	forwardToReplicas("put", "testKey", "testValue", replicas)

	// Wait for all 3 responses
	received := []string{}
	timeout := time.After(2 * time.Second)
	for i := 0; i < 3; i++ {
		select {
		case msg := <-outputChan:
			received = append(received, msg)
		case <-timeout:
			t.Fatalf("Timeout waiting for replica messages. Got %d", len(received))
		}
	}

	// Validate that 3 unique commands were received
	if len(received) != 3 {
		t.Fatalf("Expected 3 replica messages, got %d", len(received))
	}

	// Count how many had isReplica=true
	replicaCount := 0
	for _, msg := range received {
		if strings.Contains(msg, "isReplica=true") {
			replicaCount++
		}
	}

	if replicaCount != 2 {
		t.Errorf("Expected 2 isReplica=true writes, got %d", replicaCount)
	}
}
func TestHeartbeatFailure(t *testing.T) {
	// Setup ring
	ring := hashring.NewHashRing(1, 3, fnvHash)
	registeredNodes = make(map[string]bool)     // reset globals
	heartBeatFailures = make(map[string]int)

	// Start 1 healthy mock worker
	healthy := "localhost:9201"
	broken := "localhost:9202"

	ln, err := net.Listen("tcp", healthy)
	if err != nil {
		t.Fatalf("Failed to start healthy mock worker: %v", err)
	}

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				continue
			}
			go func(c net.Conn) {
				defer c.Close()
				reader := bufio.NewReader(c)
				line, _ := reader.ReadString('\n')
				if strings.TrimSpace(line) == "ping" {
					c.Write([]byte("pong\n"))
				}
			}(conn)
		}
	}()

	// Register both nodes
	ring.AddNode(healthy)
	ring.AddNode(broken)
	registeredNodes[healthy] = true
	registeredNodes[broken] = true

	startHeartBeatMonitor(ring)

	// Wait enough time for 3 failures (9s) + buffer
	time.Sleep(11 * time.Second)

	// Now check if broken node is removed
	if ring.GetNode("testkey") == broken {
		t.Errorf("Broken node was not removed after heartBeat failure")
	}

	if _, stillThere := registeredNodes[broken]; stillThere {
		t.Errorf("Broken node still in registeredNodes")
	}
}
