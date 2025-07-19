package hashring

import (
	"strconv"

	"github.com/google/btree"
)

type hashNodeItem struct {
	hash uint64	
	nodeId string
}

func (h hashNodeItem) Less(other btree.Item) bool{
	return h.hash < other.(hashNodeItem).hash
}




type ConsistentHashRing struct {
	tree *btree.BTree
	virtualNodeCount int
	replicaCount int
	hashFunc func([]byte) uint64	// func pointer
}

type Hashring interface {
	AddNode(nodeId string)
	RemoveNode(nodeId string)
	GetNode(key string) string
}

func NewHashRing(virtualNodeCount int, replicaCount int, hashFunc func([]byte) uint64) *ConsistentHashRing {
	return &ConsistentHashRing {
		tree: btree.New(16),
		virtualNodeCount: virtualNodeCount,
		replicaCount: replicaCount,
		hashFunc: hashFunc, // generic hash func can be used
	}
}

/*
Each real node (e.g., IP:Port) is assigned multiple **virtual nodes** for better key distribution.

For each real node:
1. We generate `replicaCount` number of **virtual node IDs** like "IP:Port#0", "IP:Port#1", etc.
2. Each virtual node ID is hashed using the hash function to get a **hash ID**.
3. We insert (hash ID, nodeId) pairs into the B-tree.
   - hash ID → used for placement
   - nodeId → stored so we know which real node owns that hash point

This ensures even distribution and routing for consistent hashing.
*/

func (hr ConsistentHashRing)AddNode(nodeId string) {
	// TODO: Add data replication handling after desigining replication
	for i:= 0; i<hr.virtualNodeCount ; i++ {
			virtualID := nodeId + "#" + strconv.Itoa(i)
			hash := hr.hashFunc([]byte(virtualID))
			hr.tree.ReplaceOrInsert(hashNodeItem{hash: hash, nodeId: nodeId})
	}
}

func (hr ConsistentHashRing)GetNode(key string) string {

	//0. Tree didnt form yet no nodes there, so no node
	if hr.tree.Len() == 0 {
		return ""
	}
	
	// Get hash of key 
	keyHash := hr.hashFunc([]byte(key))


	// Find node next to hash and greater than hash of key
	
	var hashOwner hashNodeItem 	// by default hash =0 , nodeId =""
	
	// func (t *btree.BTree) AscendGreaterOrEqual(pivot btree.Item, iterator btree.ItemIterator)
	/*
		for every item which is >=pivot
		the func runs the iter func(ro lambd adont know what it is)
			if return true keep going
			if false stop there
	*/
	hr.tree.AscendGreaterOrEqual( 
		hashNodeItem{hash: keyHash},
		func(item btree.Item) bool {
			hashOwner = item.(hashNodeItem)
			// basicaly typecasting item to hasNodeItem
			return false
		})
		
		// if node owns it, make the minNode as owner
		if(hashOwner.nodeId == "") {
			minNodeItem := hr.tree.Min() 
			if minNodeItem == nil { return "" }
			hashOwner = minNodeItem.(hashNodeItem)	// type cast from btree Item to hashNodeItem
		}

		return hashOwner.nodeId
}

func (hr ConsistentHashRing) RemoveNode(nodeId string) {
	// TODO : Update  replication/ migration Later
	// a node has many virtual nodes all should be removed.
	// find all virtual nodes
	for i:= 0; i<hr.virtualNodeCount ; i++ {
		// for each virtual node, get its hashid
		virtualID := nodeId + "#" + strconv.Itoa(i)
		hash := hr.hashFunc([]byte(virtualID))
		// remove the node
		hr.tree.Delete(hashNodeItem{hash: hash, nodeId: nodeId})
	}

}



func (hr ConsistentHashRing)GetReplicas(key string) []string {
	
	if hr.tree.Len() == 0 {
		return []string{}
	}

	replicas := []string{}
	seen := make(map[string]bool)
	
	// For given key find hash
	keyHash := hr.hashFunc([]byte(key))

	var hashOwner hashNodeItem 	

	hr.tree.AscendGreaterOrEqual( 
		hashNodeItem{hash: keyHash},
		func(item btree.Item) bool {
			hashOwner = item.(hashNodeItem)	// typecasting item to hasNodeItem
			nodeId := hashOwner.nodeId
			if(!seen[nodeId]) {
				replicas = append(replicas,nodeId)
				seen[nodeId] = true
			}
			return len(replicas)<hr.replicaCount
		})
	if len(replicas) < hr.replicaCount {
		hr.tree.AscendLessThan( 
			hashNodeItem{hash: keyHash},
			func(item btree.Item) bool {
				hashOwner = item.(hashNodeItem)	// typecasting item to hasNodeItem
				nodeId := hashOwner.nodeId
				if !seen[nodeId] {
					replicas = append(replicas,nodeId)
					seen[nodeId] = true
				}
				return len(replicas)<hr.replicaCount
			})
	}
	return replicas
}