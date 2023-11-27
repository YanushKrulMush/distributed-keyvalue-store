package nodes

import (
	"fmt"
	"net"
	"net/http"
	"net/rpc"
	"sync"
)

const basePort int = 12345

// KeyValue represents a key-value pair.
type KeyValue struct {
	Key   string
	Value string
}

// Node represents a node in the distributed key-value store.
type Node struct {
	ID            int
	KeyValueStore map[string]string
	mutex         sync.Mutex
}

// NewNode creates a new Node instance.
func NewNode(id int) *Node {
	return &Node{
		ID:            id,
		KeyValueStore: make(map[string]string),
	}
}

// GetKeyValue retrieves the value associated with a key.
func (n *Node) GetKeyValue(key string, value *string) error {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	val, ok := n.KeyValueStore[key]
	if !ok {
		return fmt.Errorf("Key not found")
	}

	*value = val
	return nil
}

// SetKeyValue sets the value associated with a key.
func (n *Node) SetKeyValue(kv KeyValue, reply *bool) error {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	n.KeyValueStore[kv.Key] = kv.Value
	*reply = true
	return nil
}

// StartRPCServer starts the RPC server for the node.
func (n *Node) StartRPCServer() {
	rpc.Register(n)
	rpc.HandleHTTP()

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", basePort+n.ID))
	if err != nil {
		fmt.Println("Error starting RPC server:", err)
		return
	}

	fmt.Printf("Node %d RPC server listening on port %d\n", n.ID, n.ID)
	err = http.Serve(listener, nil)
	if err != nil {
		fmt.Println("Error serving RPC server:", err)
	}
}
