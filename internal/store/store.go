package store

import (
	"DistKV/internal/cache"
	"encoding/json"
	"fmt"
	"os"
)

// TODO: Replace hardcoded "data.json" with dynamic filename per worker
const dataFile = "data.json"

var (
    kv  map[string]string
    lru *cache.LRUCache
)
func Init() {
	kv = make(map[string]string)
	load()
}

func load() {
	lru = cache.NewLRUCache(5) 
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
	lru.Put(key, value)
	kv[key] = value;
	save()
}

func Get(key string) (string, bool) { 
	if value, ok := lru.Get(key); ok {
		fmt.Println("Cache Hit")
		return value, true
	}

	if val, ok := kv[key]; ok {
		fmt.Println("Cache Miss->Scan Disk")
		lru.Put(key, kv[key])
		return val, true
	}
	// fmt.Println("Key does not exist")
	return "", false
}

func Delete(key string) bool {
	_, ok := kv[key]
	if ok {
		delete(kv, key)
		lru.Delete(key)
		save()
	}
	return ok
}

func Update(key, value string) bool {
	_, ok := kv[key]
    if ok {
        kv[key] = value
		lru.Put(key, kv[key])
        save()
    }
	return ok
}

