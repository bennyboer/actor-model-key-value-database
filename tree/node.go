package tree

import "github.com/AsynkronIT/protoactor-go/actor"

type storage map[int32]string

const CAPACITY int = 3

type Node struct {
	searchkey int32
	values    *storage
	behavior  actor.Behavior
}

func (state *Node) Receive(context actor.Context) {
	state.behavior.Receive(context)
}
