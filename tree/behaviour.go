package tree

import "github.com/AsynkronIT/protoactor-go/actor"

type Behaviour interface {
	handler(context actor.Context, node *Node)
}

// Struct to represent NodeBehavior.
type NodeBehaviour struct{}

// Struct to represent LeafBehavior
type LeafBehaviour struct{}

func createNewNode() actor.Actor {
	node := Node{nil, nil, 0, make(storage), NodeBehaviour{}}
	return &node
}

// Handles all messages for the NodeBehaviour (Actor is a Node with Leafs)
func (NodeBehaviour) handler(context actor.Context, node *Node) {

}

// Handles all messages for the LeafBehaviour (Actor is a Leaf and saves values)
func (LeafBehaviour) handler(context actor.Context, node *Node) {

}
