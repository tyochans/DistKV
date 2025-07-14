package main

import (
	"DistKV/internal/cache"
	"DistKV/internal/store"
	"bufio"
	"fmt"
	"os"
	"strings"
)

func test_lru() {

	lru := cache.NewLRUCache(2)
    lru.Put("a", "1")
    lru.Put("b", "2")
    lru.Print() // should show: b → a

    lru.Get("a")
    lru.Print() // should show: a → b

    lru.Put("c", "3")
    lru.Print() // should show: c → a (b evicted)

	valA, okA := lru.Get("a")
    fmt.Println("Get a:", valA, okA) // should return "1", true
	valB, okB := lru.Get("b")
    fmt.Println("Get b:", valB, okB) // should return "", false
	
}
func main() {

	reader:= bufio.NewReader(os.Stdin)
	fmt.Println("Welcome to Distributed KV")
	fmt.Println("Type Commands : {put, get, delete, exit, update}")
	
	store.Init()

	for {
		fmt.Print("> ")

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if (input == "exit") {
			fmt.Println("Exiting. Bye 👋")
            break
		}
		
		parts:= strings.SplitN(input, " ",3)


		if len(parts)< 2 {
			fmt.Println("Invalid Command Format")
			continue
		}
		command := parts[0]
		key :=parts[1]
		value := ""
		if len(parts) == 3{
			value = parts[2]
		}

		switch command {
		case "put":
			if value == "" {
				fmt.Println("Usage: put <key> <value>")
                continue
			}
			store.Put(key, value)
			fmt.Println("Put:", key, "->", value)
	
		case "get":
			val, ok := store.Get(key)
			if ok {
				fmt.Println("Key: ",key,  " Value:", val)
			} else {
				fmt.Println("Key not found")
			}
		case "delete":
			if store.Delete(key) {
				fmt.Println("Deleted: {",key,"}" )
			} else {
				fmt.Println("Key doesnt exist")
			}

		case "update":
			
			if store.Update(key, value) {
				fmt.Println("Updated: {",key," , ", value,"}" )
			} else {
				fmt.Println("Key doesnt exist")
			}
		default:
			fmt.Println("Unknown Command", command)

		}

	}
}
