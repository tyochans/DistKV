# DistKV
An implementation  of Distributed Key value store using GO language

## ✅ Current Status (Up to `v0.4`)

| Milestone | Status ✅ | Description                                                   |
| --------- | -------- | ------------------------------------------------------------- |
| v0.1      | ✅ Done   | Built a CLI key-value store in Go                            |
| v0.2      | ✅ Done   | Added JSON-based disk persistence + in-memory LRU cache      |
| v0.3      | ✅ Done   | Created a TCP worker server (slave) that listens on port     |
| v0.3.1    | ✅ Done   | Added a Go CLI client that connects to worker and sends commands |
| v0.4      | ✅ Done   | Added Coordinator with consistent hashing + request routing  |



### v0.4 Features

- Consistent hashing with virtual nodes (`btree`)
- Dynamic worker registration & deregistration
- Key routing from coordinator to correct worker
- TCP forwarding of `put`, `get`, `delete`, `update`
- Port passed to workers via CLI flag
