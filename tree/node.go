package tree

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"time"
)

type storage map[int32]string

const CAPACITY int = 3
const TIMEOUT time.Duration = 1000

type Node struct {
	searchkey int32
	values    *storage
	behavior  actor.Behavior
}

func (state *Node) Receive(context actor.Context) {
	state.behavior.Receive(context)
}
