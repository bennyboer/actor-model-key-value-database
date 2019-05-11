package tree

import (
	"errors"
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/ob-vss-ss19/blatt-3-sudo/messages"
	"log"
	"math"
	"sort"
)

// Creates a new node. Producer function.
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

// Utility function to forward a message to the indented child.
func (n *Node) forwardKeyedMessage(context actor.Context, key int32) error {
	var address *actor.PID

	if key <= n.searchkey {
		if len(context.Children()) > 0 {
			address = context.Children()[0]
		} else {
			return errors.New("cannot find left child of node")
		}
	} else {
		if len(context.Children()) > 1 {
			address = context.Children()[1]
		} else {
			return errors.New("cannot find right child of node")
		}
	}

	if address == nil {
		return errors.New("could not forward message to child of node")
	}

	context.Forward(address)

	return nil
}

// Behavior for leafs. Stores values and can become a node if its storage capacity is full.
func (n *Node) LeafBehavior(context actor.Context) {
	switch msg := context.Message().(type) {
	case *messages.SearchRequest:
		if val, ok := (*n.values)[msg.Key]; ok {
			var keyValue = messages.KeyValuePair{Key: msg.Key, Value: val}
			context.Respond(&messages.SearchResponse{
				Success: true,
				Entry:   &keyValue,
			})
		} else {
			context.Respond(&messages.SearchResponse{
				Success:      false,
				Entry:        nil,
				ErrorMessage: fmt.Sprintf("Could not find entry with key %d", msg.Key),
			})
		}
	case *messages.RemoveRequest:
		if val, ok := (*n.values)[msg.Key]; ok {
			removed := messages.KeyValuePair{Key: msg.Key, Value: val}
			delete(*n.values, msg.Key)
			context.Respond(&messages.RemoveResponse{
				Success:     true,
				RemovedPair: &removed,
			})
		} else {
			context.Respond(&messages.RemoveResponse{
				Success:      false,
				RemovedPair:  nil,
				ErrorMessage: fmt.Sprintf("Could not find key %d to delete", msg.Key),
			})
		}
	case *messages.InsertRequest:
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
			n.values = nil
			n.behavior.Become(n.NodeBehavior)
			log.Printf("[Node] Leaf is now a node")
		}

		context.Respond(&messages.InsertResponse{
			Success: true,
		})
	case *messages.TraverseRequest:
		log.Println("Traverse LEAF")
		var pairs = make([]*messages.KeyValuePair, 0, len(*n.values))

		log.Println("Iterating over pairs")
		for k, v := range *(n.values) {
			pairs = append(pairs, &messages.KeyValuePair{
				Key:   k,
				Value: v,
			})
		}

		log.Println("Sort pairs")
		sort.Slice(pairs, func(i, j int) bool {
			return pairs[i].Key < pairs[j].Key
		})

		log.Println("Respond TraverseResponse")
		context.Respond(&messages.TraverseResponse{
			Success: true,
			Pairs:   pairs,
		})
	}
}

// Behavior for the node. Acts as a parent to two leafs and forwards indented messages.
func (n *Node) NodeBehavior(context actor.Context) {
	switch msg := context.Message().(type) {
	case *messages.SearchRequest:
		err := n.forwardKeyedMessage(context, msg.Key)
		if err != nil {
			context.Respond(&messages.SearchResponse{
				Success:      false,
				Entry:        nil,
				ErrorMessage: err.Error(),
			})
		}
	case *messages.InsertRequest:
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
		err := n.forwardKeyedMessage(context, msg.Key)
		if err != nil {
			context.Respond(&messages.RemoveResponse{
				Success:      false,
				RemovedPair:  nil,
				ErrorMessage: err.Error(),
			})
		}
	case *messages.DeleteTreeRequest:
		for _, child := range context.Children() {
			context.Send(child, &messages.DeleteTreeRequest{})
			child.Poison()
		}
	case *messages.TraverseRequest:
		log.Println("Traverse NODE")
		pairs := make([]*messages.KeyValuePair, 0)
		for _, child := range context.Children() {
			result, _ := context.RequestFuture(child, &messages.TraverseRequest{}, TIMEOUT).Result()
			response, ok := result.(*messages.TraverseResponse)

			if ok {
				pairs = append(pairs, response.Pairs...)
			} else {
				log.Println("Child node responded with incompatible type. Expected TraverseResponse")
				context.Respond(&messages.TraverseResponse{
					Success:      false,
					Pairs:        nil,
					ErrorMessage: "FATAL: Child node responded with incompatible type",
				})
			}
		}

		context.Respond(&messages.TraverseResponse{
			Success: true,
			Pairs:   pairs,
		})
	}
}
