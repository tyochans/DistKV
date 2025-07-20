package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"testing"
	"time"
)

func TestAutoRegister(t *testing.T) {
	regMsgChan := make(chan string, 1)
	errChan := make(chan error, 1)

	// Step 1: Start fake coordinator (run safely without t.Fatalf)
	go func() {
		ln, err := net.Listen("tcp", "localhost:9310")
		if err != nil {
			errChan <- fmt.Errorf("Coordinator listen failed: %v", err)
			return
		}
		defer ln.Close()

		conn, err := ln.Accept()
		if err != nil {
			errChan <- fmt.Errorf("Coordinator accept failed: %v", err)
			return
		}
		defer conn.Close()

		msg, _ := bufio.NewReader(conn).ReadString('\n')
		regMsgChan <- strings.TrimSpace(msg)
	}()

	// Step 2: Start worker that will auto-register
	startAutoRegister("localhost:9310", "localhost:9400")

	// Step 3: Wait for either success, error, or timeout
	select {
	case err := <-errChan:
		t.Fatal(err)
	case msg := <-regMsgChan:
		if !strings.HasPrefix(msg, "register localhost:9400") {
			t.Errorf("Unexpected register message: %s", msg)
		}
	case <-time.After(5 * time.Second):
		t.Fatalf("Timed out waiting for auto-register")
	}
}
