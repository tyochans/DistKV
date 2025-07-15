package main

import (
	"DistKV/internal/store"
	"bufio"
	"fmt"
	"net"
	"strings"
)

func main() {

	listener, err:= net.Listen("tcp", ":9000")
	// why is it not integer and ":9000" with a colon?

	if err!=nil {
		panic(err) // why panic?
	}	
	fmt.Println("Worker Server started on  port 9000")
	
	for {
		connection, err := listener.Accept()
		if err!=nil {
			continue
		}	
		go handleConnection(connection)

		
	}
	
}

func handleConnection(connection net.Conn) {
//WHy this is not in Caps?
	defer connection.Close()
	fmt.Println("Client Connected")
	reader := bufio.NewReader(connection)
	// So this is not for JSON right?
	store.Init()
	for {
		fmt.Print("> ")
		input, err := reader.ReadString('\n')

		if err !=nil {
			fmt.Println("Connection Error")
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
		default:
			connection.Write([]byte("Unknown Command: " + command +"\n"))

		}
	}
}