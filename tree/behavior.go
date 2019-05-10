package tree

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/ob-vss-ss19/blatt-3-sudo/messages"
	"log"
	"math"
	"sort"
)

func NewNode(capacity int, values *Storage) actor.Actor {
	act := &Node{
		capacity:  capacity,
		values:    values,
		searchkey: math.MinInt32,
		behavior:  actor.NewBehavior(),
	}

	act.behavior.Become(act.LeafBehavior)

	return act
}

func (n *Node) forwardKeyedMessage(context actor.Context, key int32) {
	var address *actor.PID = nil
	if n.searchkey <= key {
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

func (n *Node) LeafBehavior(context actor.Context) {
	switch msg := context.Message().(type) {
	case *messages.SearchRequest:
		if val, ok := (*n.values)[msg.Key]; ok {
			var keyValue = messages.KeyValuePair{Key: msg.Key, Value: val}
			context.Respond(&messages.SearchResponse{Success: true, Entry: &keyValue})
		} else {
			context.Respond(&messages.SearchResponse{Success: false, Entry: nil})
		}
	case *messages.RemoveRequest:
		if val, ok := (*n.values)[msg.Key]; ok {
			removed := messages.KeyValuePair{Key: msg.Key, Value: val}
			delete(*n.values, msg.Key)
			context.Respond(&messages.RemoveResponse{Success: true, RemovedPair: &removed})
		} else {
			context.Respond(&messages.RemoveResponse{Success: false, RemovedPair: nil})
		}
	case *messages.InsertRequest:
		log.Println("Insert into leaf")
		if len(*n.values) <= n.capacity {
			(*n.values)[msg.Entry.Key] = msg.Entry.Value
			log.Printf("[Node] Inserted %d with %v", msg.Entry.Key, msg.Entry.Value)
		} else {
			// Sort key values pairs by keys
			(*n.values)[msg.Entry.Key] = msg.Entry.Value
			var pairs = make([]messages.KeyValuePair, 0, n.capacity+1)

			for k, v := range *n.values {
				pairs = append(pairs, messages.KeyValuePair{
					Key:   k,
					Value: v,
				})
			}

			sort.Slice(pairs, func(i, j int) bool {
				return pairs[i].Key < pairs[j].Key
			})

			// Split storage in two
			midIndex := len(pairs) / 2
			isEven := len(pairs)%2 == 0
			if !isEven {
				midIndex++
			}
			n.searchkey = pairs[midIndex].Key

			leftValues := make(Storage, n.capacity+1)
			rightValues := make(Storage, n.capacity+1)

			for i := range pairs {
				pair := pairs[i]

				if i <= midIndex {
					leftValues[pair.Key] = pair.Value
				} else {
					rightValues[pair.Key] = pair.Value
				}
			}

			// Spawn both children
			context.Spawn(actor.PropsFromProducer(func() actor.Actor {
				return NewNode(n.capacity, &leftValues)
			}))
			context.Spawn(actor.PropsFromProducer(func() actor.Actor {
				return NewNode(n.capacity, &rightValues)
			}))

			// Leaf is now a node
			(*n).values = nil
			n.behavior.Become(n.NodeBehavior)
			log.Printf("[Node] Leaf is now a node")
		}
		context.Respond(&messages.InsertResponse{Success: true})
	case *messages.TraverseRequest:
		log.Println("Traverse leaf")
		var pairs = make([]*messages.KeyValuePair, 0, len(*n.values))

		for k, v := range *(n.values) {
			pairs = append(pairs, &messages.KeyValuePair{
				Key:   k,
				Value: v,
			})
		}

		sort.Slice(pairs, func(i, j int) bool {
			return pairs[i].Key < pairs[j].Key
		})

		context.Respond(&messages.TraverseResponse{
			Pairs: pairs,
		})
	}
}

func (n *Node) NodeBehavior(context actor.Context) {
	switch msg := context.Message().(type) {
	case *messages.SearchRequest:
		n.forwardKeyedMessage(context, msg.Key)

	case *messages.InsertRequest:
		println("Insert into node")
		var address *actor.PID
		if n.searchkey < msg.Entry.Key {
			if len(context.Children()) > 1 {
				address = context.Children()[1]
			} else {
				values := make(Storage, n.capacity+1)
				address = context.Spawn(actor.PropsFromProducer(func() actor.Actor {
					return NewNode(n.capacity, &values)
				}))
			}
		} else {
			if len(context.Children()) > 0 {
				address = context.Children()[0]
			} else {
				values := make(Storage, n.capacity+1)
				address = context.Spawn(actor.PropsFromProducer(func() actor.Actor {
					return NewNode(n.capacity, &values)
				}))
			}
		}
		context.Forward(address)
	case *messages.RemoveRequest:
		n.forwardKeyedMessage(context, msg.Key)
	case *messages.DeleteTreeRequest:
		for _, child := range context.Children() {
			context.Send(child, &messages.DeleteTreeRequest{})
			child.Poison()
		}
	case *messages.TraverseRequest:
		log.Println("Traverse Node")
		pairs := make([]*messages.KeyValuePair, 0)
		for _, child := range context.Children() {
			log.Println("Traverse next Node")
			result, _ := context.RequestFuture(child, &messages.TraverseRequest{}, TIMEOUT).Result()
			response, err := result.(*messages.TraverseResponse)

			if err {
				pairs = append(pairs, response.Pairs...)
			}
		}

		context.Respond(&messages.TraverseResponse{Pairs: pairs})
	}
}
