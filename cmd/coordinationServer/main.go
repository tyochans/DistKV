package main

import (
	"DistKV/cmd/coordinationServer/hashring"
	"bufio"
	"fmt"
	"hash/fnv"
	"net"
	"strings"
)


func fnvHash(data []byte)uint64 {
	hasher := fnv.New64a()
	hasher.Write(data)
	return hasher.Sum64()
}



func HandleClient(conn net.Conn, ring *hashring.ConsistentHashRing) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	for {
		conn.Write([]byte("Waiting for Command\n"))
		line, err := reader.ReadString('\n')

		if err != nil {
			fmt.Println("Client Disconnected")
			return
		}

		line = strings.TrimSpace(line)
		parts := strings.Split(line, " ")
	//TODO : Add some timer so Coordinator doesnt wait too long for the command
		if len(parts)<2 {
			conn.Write([]byte("Invalid command format"))
				continue
		}
		command := parts[0]
		switch (command) {
			case "register":
				nodeId := parts[1]
				ring.AddNode(nodeId)
				conn.Write([]byte("Registered " + nodeId + "\n"))
				
			case "deregister":
				nodeId := parts[1]
				ring.RemoveNode(nodeId)
				conn.Write([]byte("Deregistered " + nodeId + "\n"))
				
//TODO : do something such that only put get update delete are forwarded reducing uneccessary communications
			default :
				// Extract key and lookup worker
				key := parts[1]
				workerAddr := ring.GetNode(key)
				if workerAddr == "" {
					conn.Write([]byte("No worker available\n"))
					continue
				}

				// Forward command to the responsible worker
				forwardAndRelay(conn, workerAddr, line)
		}
	}
}
func forwardAndRelay(clientConn net.Conn, workerAddr, line string) {
	workerConn, err := net.Dial("tcp", workerAddr)
	
	if err != nil {
		clientConn.Write([]byte("Failed to connect to worker\n"))
		return
	}

	defer workerConn.Close()
	_, err = workerConn.Write([]byte(line + "\n"))

	if err != nil {
		clientConn.Write([]byte("Failed to forward to worker\n"))

	}
	response, _ := bufio.NewReader(workerConn).ReadString('\n')
	clientConn.Write([]byte("Worker " + workerAddr + ": " + response))
}
func main() {
	ring := hashring.NewHashRing(3, fnvHash)
	listener, err := net.Listen("tcp", ":9000")
	if err != nil {
		fmt.Println("Error Starting COordinator",err)
		return 
	}
	defer listener.Close()
	fmt.Println("Coordinator is listening on port 9000")
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error Accepting client")
			continue
		}
		fmt.Println("Client connected:", conn.RemoteAddr())
		go HandleClient(conn, ring)
	}
}
