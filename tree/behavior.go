package tree

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/ob-vss-ss19/blatt-3-sudo/messages"
	"log"
	"math"
	"sort"
)

func NewNode(capacity int) actor.Actor {
	storage := make(storage, capacity+1)
	act := &Node{
		capacity:  capacity,
		values:    &storage,
		searchkey: math.MinInt32,
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
	case *messages.SearchRequest:
		if val, ok := (*node.values)[msg.Key]; ok {
			var keyValue = messages.KeyValuePair{Key: msg.Key, Value: val}
			context.Respond(&messages.SearchResponse{Success: true, Entry: &keyValue})
		} else {
			context.Respond(&messages.SearchResponse{Success: false, Entry: nil})
		}
	case *messages.RemoveRequest:
		if val, ok := (*node.values)[msg.Key]; ok {
			removed := messages.KeyValuePair{Key: msg.Key, Value: val}
			delete(*node.values, msg.Key)
			context.Respond(&messages.RemoveResponse{Success: true, RemovedPair: &removed})
		} else {
			context.Respond(&messages.RemoveResponse{Success: false, RemovedPair: nil})
		}
	case *messages.InsertRequest:
		if len(*node.values) <= node.capacity {
			(*node.values)[msg.Entry.Key] = msg.Entry.Value
			log.Printf("[Node] Inserted %d with %v", msg.Entry.Key, msg.Entry.Value)
		} else {
			// Spawn both children
			context.Spawn(actor.PropsFromProducer(func() actor.Actor {
				return NewNode(node.capacity)
			}))
			context.Spawn(actor.PropsFromProducer(func() actor.Actor {
				return NewNode(node.capacity)
			}))

			(*node.values)[msg.Entry.Key] = msg.Entry.Value
			var keys = make([]int32, 0, node.capacity+1)

			for k := range *node.values {
				keys = append(keys, k)
			}

			sort.Slice(keys, func(i, j int) bool {
				return keys[i] < keys[j]
			})
			midIndex := len(keys) / 2
			isEven := len(keys)%2 == 0
			if !isEven {
				midIndex++
			}
			node.searchkey = keys[midIndex]

			for i, k := range keys {
				var entry = messages.KeyValuePair{Key: k, Value: (*node.values)[k]}
				var index = 0
				if i > midIndex {
					index = 1
				}
				context.Send(context.Children()[index], &messages.InsertRequest{Entry: &entry, TreeId: msg.TreeId})
			}

			// Leaf is now a node
			(*node).values = nil
			node.behavior.Become(node.NodeBehavior)
			log.Printf("[Node] Leaf is now a node")
		}
		context.Respond(&messages.InsertResponse{Success: true})
	case *messages.TraverseRequest:
		var pairs = make([]*messages.KeyValuePair, 0, 3)
		var keys = make([]int32, 0, 3)

		for k := range *(node.values) {
			keys = append(keys, k)
		}

		sort.Slice(keys, func(i, j int) bool {
			return keys[i] < keys[j]
		})

		for _, i := range keys {
			pairs = append(pairs, &messages.KeyValuePair{Key: i, Value: (*node.values)[i]})
		}

		context.Respond(&messages.TraverseResponse{Pairs: pairs})
	}
}

func (node *Node) NodeBehavior(context actor.Context) {
	switch msg := context.Message().(type) {
	case *messages.SearchRequest:
		node.forwardKeyedMessage(context, msg.Key)

	case *messages.InsertRequest:
		var address *actor.PID
		if node.searchkey < msg.Entry.Key {
			if len(context.Children()) > 1 {
				address = context.Children()[1]
			} else {
				address = context.Spawn(actor.PropsFromProducer(func() actor.Actor {
					return NewNode(node.capacity)
				}))
			}
		} else {
			if len(context.Children()) > 0 {
				address = context.Children()[0]
			} else {
				address = context.Spawn(actor.PropsFromProducer(func() actor.Actor {
					return NewNode(node.capacity)
				}))
			}
		}
		context.Forward(address)
	case *messages.RemoveRequest:
		node.forwardKeyedMessage(context, msg.Key)
	case *messages.DeleteTreeRequest:
		for _, child := range context.Children() {
			context.Send(child, &messages.DeleteTreeRequest{})
			child.Poison()
		}
	case *messages.TraverseRequest:
		var pairs = make([]*messages.KeyValuePair, 0, node.capacity)
		for _, child := range context.Children() {
			result, _ := context.RequestFuture(child, &messages.TraverseRequest{}, TIMEOUT).Result()
			response, err := result.(*messages.TraverseResponse)

			if err {
				pairs = append(pairs, response.Pairs...)
			}
		}

		context.Respond(&messages.TraverseResponse{Pairs: pairs})
	}
}
