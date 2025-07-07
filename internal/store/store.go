package store

import (
	"fmt"
	"encoding/json"
	"os"
)

const dataFile = "data.json"

var kv map[string]string

func init() {
	kv = make(map[string]string)
	load()
}

func load() {
	bytes, err := os.ReadFile(dataFile)
	if err!=nil {
		fmt.Println("No existing store found, creating new store")
		return
	}
	err = json.Unmarshal(bytes, &kv)
	if err!= nil {
		fmt.Println("Error reading store : ", err)
		return
	} else {
		fmt.Println("Loaded Existing Store")
	}

}

func save() {
	bytes, err := json.MarshalIndent(kv, "", " ")
	if err!=nil {
		fmt.Println("Failed to Serialize store")
		return
	}
	err = os.WriteFile(dataFile, bytes, 0644)
	if err!= nil {
		fmt.Println("Failed to save store : ", err)
		return
	} else {
		fmt.Println(" Store Saved")
	}
}

func Put(key, value string) {
	kv[key] = value;
	save()
}

func Get(key string) (string, bool) { 
	val, ok := kv[key]
	return val, ok
}

func Delete(key string) bool {
	_, ok := kv[key]
	if ok {
		delete(kv, key)
		save()
	}
	return ok
}

func Update(key, value string) bool {
	_, ok := kv[key]
    if ok {
        kv[key] = value
        save()
    }
	return ok
}

/*
FIle	"DistKV/internal/store.go"

1. const dataFile = "data.json" is this private and cant be seen in main

2.
var kv map[string]string

func init() {
	kv = make(map[string]string)
	load()
}
Is is like global variable but why are we writng it 2 times
cant we write just var kv = make(map[string]string) in a single line

3. 
func load() {
	bytes, err := os.ReadFile(dataFile)
	if err!=nil {
		fmt.Println("No existing store found, creating new store")
		return
	}
WHat is nil? is it similar to NULL 

4. func save() {
	bytes, err := json.MarshalIndent(kv, "", " ")

	I forgot what this does

5.
func Put(key, value string) {

here the paramters are written ina  weird way does this mean botha re of string type
is below valid paramlist
a,b int, c,d string, e,f bool

6.
	_, ok := kv[key]
	if ok {
		delete(kv, key)
Doesnt delete return any thing why we are checking firtst and then deelting cant we directly write
 ok = delete(kv, key)
 
*/