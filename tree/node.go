// The tree package provides the behavior and structs for the tree components
package tree

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"time"
)

// Type to simplify the usage of the storage
type Storage map[int32]string

const TIMEOUT = time.Second * 5

// Struct that contains all necessary variables for a tree element
type Node struct {
	searchkey int32          // Max left searchkey
	values    *Storage       // Value storage
	behavior  actor.Behavior // The current behavior of the node
	capacity  int            // The nodes maximum capacity
}

// Implements the actor interface to receive messages
func (n *Node) Receive(context actor.Context) {
	n.behavior.Receive(context)
}
