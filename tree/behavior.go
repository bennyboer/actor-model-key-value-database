package tree

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/ob-vss-ss19/blatt-3-sudo/messages"
)

func newNode() actor.Actor {
	act := &Node{
		behavior: actor.NewBehavior(),
	}
	act.behavior.Become(act.LeafBehavior)
	return act
}

func (node *Node) LeafBehavior(context actor.Context) {
	switch msg := context.Message().(type) {
	case messages.SearchRequest:
		// TODO LeafBehavior
	}
}

func (node *Node) forwardKeyedMessage(context actor.Context, key int32) {
	var address *actor.PID = nil
	if node.searchkey <= key {
		if len(context.Children()) > 1 {
			address = context.Children()[1]
		}
	} else {
		if len(context.Children()) > 0 {
			address = context.Children()[0]
		}
	}
	if address != nil {
		context.Send(address, context.Message())
	}
}

func (node *Node) NodeBehavior(context actor.Context) {
	switch msg := context.Message().(type) {
	case messages.SearchRequest:
		node.forwardKeyedMessage(context, msg.Key)
	case messages.InsertRequest:
		var address *actor.PID
		if node.searchkey <= msg.Entry.Key {
			if len(context.Children()) > 1 {
				address = context.Children()[1]
			} else {
				address = context.Spawn(actor.PropsFromProducer(newNode))
			}
		} else {
			if len(context.Children()) > 0 {
				address = context.Children()[0]
			} else {
				address = context.Spawn(actor.PropsFromProducer(newNode))
			}
		}
		context.Send(address, context.Message())
	case messages.RemoveRequest:
		node.forwardKeyedMessage(context, msg.Key)
		// TODO DeleteTree
	}
}
