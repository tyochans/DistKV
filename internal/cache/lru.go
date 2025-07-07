package cache

import (
	"fmt"
)

type LRUCache struct {
	capacity int
	cache    map[string]*Node
	head     *Node
	tail     *Node
}

type Node struct {
    key   string
    value string
    prev  *Node
    next  *Node
}

func NewLRUCache(capacity int) *LRUCache {
	return &LRUCache {
		capacity:capacity,
		cache: make(map[string]*Node),
		head: nil,
		tail: nil,
	}

}

func (lru *LRUCache) Print() {
	curr := lru.head

	for curr!= nil {
		fmt.Print("[%s:%s] -> ", curr.key, curr.value)
		curr = curr.next
	}
	fmt.Println(nil)
}

func (lru *LRUCache) Put(key, value string) {
	
	// checked if key already there
	// update and move it to front
	if node, exists := lru.cache[key];exists {
		node.value = value
		lru.moveToFront(node)
		return
	}
	// if not there need to make a new one
	// make one
	newNode := &Node{
		key: key,
		value: value,
	} 
	// add to cache
	//move it to front
	lru.cache[key] = newNode
	lru.addToFront(newNode)

	// if overflow remove last in list
	if(len(lru.cache) >lru.capacity) {
		lru.removeTail()
	}
}

func (lru *LRUCache) addToFront(node *Node) {
	node.prev = nil
	node.next = lru.head
	if lru.head!= nil {
		lru.head.prev = node
	}
	lru.head = node
}
func (lru *LRUCache) moveToFront(node *Node) {
    if node == lru.head {
        return // already at front
    }	

	if node.prev != nil {
		node.prev.next = node.next
	}

	if node.next != nil {
		node.next.prev = node.prev
	}

	if node == lru.tail {
		lru.tail = node.prev
	}

	node.prev = nil
	node.next = lru.head
	if lru.head!= nil {
		lru.head.prev = node
	}
	lru.head = node
}
func (lru *LRUCache) removeTail(){
	if lru.tail!=nil {
		temp := lru.tail
		if temp.prev != nil {
			temp.prev.next = nil
			lru.tail = temp.prev
		} else {
			lru.head = nil
			lru.tail = nil
		}
		delete(lru.cache, temp.key)
	}
}


func (lru *LRUCache) Get(key string) (string, bool){
	node, exists := lru.cache[key]
	if !exists {
		return "", false
	}
	lru.moveToFront(node)
	return node.value, true
}