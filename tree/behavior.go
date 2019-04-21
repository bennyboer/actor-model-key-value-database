package tree

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/ob-vss-ss19/blatt-3-sudo/messages"
	"math"
)

func newNode() actor.Actor {
	storage := make(storage, 3)
	act := &Node{
		values:    &storage,
		searchkey: math.MinInt32,
		left:      nil,
		right:     nil,
		behavior:  actor.NewBehavior(),
	}
	act.behavior.Become(act.LeafBehavior)
	return act
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
		context.Forward(address)
	}
}

func (node *Node) LeafBehavior(context actor.Context) {
	switch msg := context.Message().(type) {
	case messages.SearchRequest:
		if val, ok := (*node.values)[msg.Key]; ok {
			var keyValue = messages.KeyValuePair{Key: msg.Key, Value: val}
			context.Respond(messages.SearchResponse{Success: true, Entry: &keyValue})
		} else {
			context.Respond(messages.SearchResponse{Success: false, Entry: nil})
		}
	case messages.RemoveRequest:
		if val, ok := (*node.values)[msg.Key]; ok {
			removed := messages.KeyValuePair{Key: msg.Key, Value: val}
			delete(*node.values, msg.Key)
			context.Respond(messages.RemoveResponse{Success: true, RemovedPair: &removed})
		} else {
			context.Respond(messages.RemoveResponse{Success: false, RemovedPair: nil})
		}
	case messages.InsertRequest:
		if len(*node.values) < CAPACITY {
			(*node.values)[msg.Entry.Key] = msg.Entry.Value
		} else {
			context.Spawn(actor.PropsFromProducer(newNode))
			node.searchkey = math.MinInt32
			for k, v := range *node.values {
				var entry = messages.KeyValuePair{Key: k, Value: v}

				if k > node.searchkey {
					node.searchkey = k
				}

				context.Send(context.Children()[0], messages.InsertRequest{Entry: &entry, TreeId: msg.TreeId})
			}
			(*node).values = nil
			// Leaf is now a node
			node.behavior.Become(node.NodeBehavior)
		}
		context.Respond(messages.InsertResponse{Success: true})
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
		context.Forward(address)
	case messages.RemoveRequest:
		node.forwardKeyedMessage(context, msg.Key)
		// TODO DeleteTree
	}
}
