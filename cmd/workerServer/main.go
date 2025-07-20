package main

import (
	"DistKV/internal/store"
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

func main() {

	port := flag.String("port", "9000", "Port for the worker server to listen on")
	flag.Parse()

	address := ":" + *port

	coordinatorAddr := flag.String("coordinator", "localhost:9000", "Address of the coordinator")
	workerAddr := "localhost:" + *port
	
	startAutoRegister(*coordinatorAddr, workerAddr)

	listener, err:= net.Listen("tcp", address)

	if err!=nil {
		fmt.Println("Error starting worker on port", *port, ":", err)
		os.Exit(1) 
	}	
	fmt.Println("Worker Server started on  port " + *port)
	
	for {
		connection, err := listener.Accept()
		if err!=nil {
			continue
		}	
		go handleConnection(connection)

		
	}
	
}

func startAutoRegister(coordinatorAddr, workerAddr string) {
	go func() {
		for {
			conn, err := net.Dial("tcp", coordinatorAddr)
			if err!=nil {
				fmt.Println("Failed to connect to coordinator:", err)
				time.Sleep(5 * time.Second)
				continue
			}
			fmt.Fprintf(conn, "register %s\n", workerAddr)

			response,_ := bufio.NewReader(conn).ReadString('\n')
			fmt.Println("Coordinator response", strings.TrimSpace(response))

			conn.Close()
			return
		}
	}()
}
func handleConnection(connection net.Conn) {
	defer connection.Close()
	fmt.Println("Coordinator Connected")
	reader := bufio.NewReader(connection)
	store.Init()
	for {
		input, err := reader.ReadString('\n')

		if err !=nil {
			fmt.Println("Connection Error")
			return
		}
		input = strings.TrimSpace(input)
		
		if (input == "exit") {
			connection.Write([]byte("Exiting. Bye "))
            break
		}
		
		parts:= strings.SplitN(input, " ",3)


		if len(parts)< 2 {
			connection.Write([]byte("Invalid Command Format"))
			continue
		}
		command := strings.ToLower(parts[0])
		fmt.Print("Command : "+ command)
		key :=parts[1]
		fmt.Print(" ",key)
		value := ""
		if len(parts) == 3{
			value = parts[2]
			fmt.Print(" , ",value)		
		}
		fmt.Println("")
		switch command {
		case "put":
			if value == "" {
				connection.Write([]byte("Usage: put <key> <value>\n"))
                continue
			}
			store.Put(key, value)
			connection.Write([]byte("Put: " + key + " -> " + value+"\n"))
	
		case "get":
			val, ok := store.Get(key)
			if ok {
				connection.Write([]byte("Key: " + key + " Value: " + val + "\n"))
			} else {
				connection.Write([]byte("Key not found" + "\n"))
			}
		case "delete":
			if store.Delete(key) {
				connection.Write([]byte("Deleted: {" + key + "}\n"))
			} else {
				connection.Write([]byte("Key doesnt exist\n"))
			}

		case "update":
			
			if store.Update(key, value) {
				connection.Write([]byte("Updated: {" + key + " , " + value + "}\n"))
			} else {
				connection.Write([]byte("Key doesnt exist\n"))
			}
		case "ping":
				connection.Write([]byte("pong\n"))
		default:
			connection.Write([]byte("Unknown Command: " + command +"\n"))

		}
	}
}