package main

import (
	"DistKV/cmd/coordinationServer/hashring"
	"bufio"
	"flag"
	"fmt"
	"hash/fnv"
	"net"
	"strings"
	"time"
)


func fnvHash(data []byte)uint64 {
	hasher := fnv.New64a()
	hasher.Write(data)
	return hasher.Sum64()
}

func forwardToReplicas(cmd string, key string, value string, replicas []string){

	for i, addr := range replicas {
		go func(addr string, isReplica bool) {
			conn, err := net.Dial("tcp", addr)
			if err != nil {
				fmt.Printf("Failed to connect to %s: %v\n", addr, err)
				return
			}
			defer conn.Close()

			// Construct command string
			fullCmd := fmt.Sprintf("%s %s %s", cmd, key, value)
			if isReplica {
				fullCmd += " isReplica=true"
			}
			fullCmd += "\n"
			_, err = conn.Write([]byte(fullCmd))
			if err != nil {
				fmt.Printf("Failed to send to %s: %v\n", addr, err)
			}
		}(addr, i>0)
	}
}

func HandleClient(conn net.Conn, ring *hashring.ConsistentHashRing) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	for {
		conn.Write([]byte("Waiting for Command\n"))

		conn.SetReadDeadline(time.Now().Add(10 * time.Second))
		line, err := reader.ReadString('\n')

		if err != nil {
			fmt.Println("Client Disconnected")
			return
		}

		line = strings.TrimSpace(line)
		parts := strings.Split(line, " ")
		if len(parts)<2 {
			conn.Write([]byte("Invalid command format"))
				continue
		}
		command := parts[0]
		switch (command) {
			case "register":
				nodeId := parts[1]
				ring.AddNode(nodeId)
				registeredNodes[nodeId] = true;
				conn.Write([]byte("Registered " + nodeId + "\n"))
				
			case "deregister":
				nodeId := parts[1]
				ring.RemoveNode(nodeId)
				conn.Write([]byte("Deregistered " + nodeId + "\n"))
				
			default :
				// Extract key and lookup worker
				valid := map[string]bool{"put": true, "get": true, "update": true, "delete": true}
				if !valid[command] {
					conn.Write([]byte("Invalid command\n"))
					continue
				}
				key := parts[1]
				if command == "put" && len(parts) >= 3 {
					value := parts[2]
					replicas := ring.GetReplicas(key)
					forwardToReplicas("put", key, value, replicas)
					conn.Write([]byte("PUT forwarded to replicas\n"))
	 			} else {
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

	port := flag.String("port", "9000", "Port for the Coordination server to listen on")
	virtualNodeCount := flag.Int("Virtual Nodes", 3, "No of virutal nodes per Server")
	replicaCount := flag.Int("Replica Count", 3, "No of replicas of a worker (including primary)")

	address := ":" + *port
	ring := hashring.NewHashRing(*virtualNodeCount, *replicaCount, fnvHash)
	startHeartBeatMonitor(ring)
	listener, err := net.Listen("tcp",address)
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

var registeredNodes = make(map[string]bool)
var heartBeatFailures = make(map[string]int)


func startHeartBeatMonitor(ring *hashring.ConsistentHashRing){
	ticker := time.NewTicker(3* time.Second)
	go func() {
		for range ticker.C {
			for nodeId := range registeredNodes {
				go func(nodeId string) {
					conn, err := net.DialTimeout("tcp",nodeId,time.Second)
					if err!= nil {
						heartBeatFailures[nodeId]++
					}else {
						defer conn.Close()
						conn.SetWriteDeadline(time.Now().Add(1*time.Second))
						fmt.Fprintln(conn, "ping")

						reply, err:= bufio.NewReader(conn).ReadString('\n')
						if err!=nil || strings.TrimSpace(reply)!="pong" {
							heartBeatFailures[nodeId]++
						} else{
							heartBeatFailures[nodeId]=0
						}	
					}
					if heartBeatFailures[nodeId] >= 3 {
						fmt.Println("Node Dead: ",nodeId)
						ring.RemoveNode(nodeId)
						delete(registeredNodes,nodeId)
						delete(heartBeatFailures,nodeId)
					}
				}(nodeId)
			}
		}
	}()
}