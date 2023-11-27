package nodes

import "net/rpc"

// NodeClient represents a client for interacting with a node in the distributed key-value store.
type NodeClient struct {
	client *rpc.Client
}

// NewNodeClient creates a new NodeClient instance.
func NewNodeClient(address string) (*NodeClient, error) {
	client, err := rpc.DialHTTP("tcp", address)
	if err != nil {
		return nil, err
	}
	return &NodeClient{client: client}, nil
}

// GetKeyValue retrieves the value associated with a key from the node.
func (c *NodeClient) GetKeyValue(key string) (string, error) {
	var value string
	err := c.client.Call("Node.GetKeyValue", key, &value)
	return value, err
}

// SetKeyValue sets the value associated with a key on the node.
func (c *NodeClient) SetKeyValue(keyValue KeyValue) error {
	var reply bool
	err := c.client.Call("Node.SetKeyValue", keyValue, &reply)
	return err
}
