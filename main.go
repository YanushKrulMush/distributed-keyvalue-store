package main

import (
	"distributed-keyvalue-store/nodes"
	"fmt"
)

func main() {
	node := nodes.NewNode(1)
	go node.StartRPCServer()

	// Replace these addresses with the actual addresses of your nodes.
	node1Address := "localhost:12346"
	// node2Address := "localhost:12346"

	client1, err := nodes.NewNodeClient(node1Address)
	if err != nil {
		fmt.Println("Error creating client for node 1:", err)
		return
	}

	// client2, err := nodes.NewNodeClient(node2Address)
	// if err != nil {
	// 	fmt.Println("Error creating client for node 2:", err)
	// 	return
	// }

	// Use the clients to make RPC calls to different nodes
	// Example: Set a key-value pair on node 1
	err = client1.SetKeyValue(nodes.KeyValue{Key: "exampleKey", Value: "exampleValue"})
	if err != nil {
		fmt.Println("Error setting key-value pair on node 1:", err)
		return
	}

	// Example: Get the value associated with a key from node 2
	value, err := client1.GetKeyValue("exampleKey")
	if err != nil {
		fmt.Println("Error getting value from node 2:", err)
		return
	}

	fmt.Printf("Value for key 'exampleKey' from node 2: %s\n", value)

	// Keep the node running
	select {}
}
