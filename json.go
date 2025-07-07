package main

import ( 
	"encoding/json"
	"fmt"
	"os"
)

func main() {
	store:= map[string]string {
		"foo" : "bar",
		"baz" : "qux",
	}

	jsonBytes, _ :=json.MarshalIndent(store, "", " ")
	os.WriteFile("store.json", jsonBytes, 0644)

	fmt.Println("Save to store.json")


	loadedBytes , _ := os.ReadFile("store.json")
	var loaded map[string]string
	json.Unmarshal(loadedBytes, &loaded)
	fmt.Println("Loaded\n", loaded)
}