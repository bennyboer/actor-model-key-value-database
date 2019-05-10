package tree

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"time"
)

type Storage map[int32]string

const TIMEOUT time.Duration = 1000

type Node struct {
	searchkey int32
	values    *Storage
	behavior  actor.Behavior
	capacity  int
}

func (n *Node) Receive(context actor.Context) {
	n.behavior.Receive(context)
}
