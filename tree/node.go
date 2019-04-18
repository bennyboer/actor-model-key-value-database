package tree

import "github.com/AsynkronIT/protoactor-go/actor"

type storage map[int]string

type Node struct {
	left      *Node
	right     *Node
	key       int
	values    storage
	behaviour Behaviour
}

func (n *Node) Receive(context actor.Context) {
	n.behaviour.handler(context, n)
}
