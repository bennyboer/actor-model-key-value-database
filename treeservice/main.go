package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/AsynkronIT/goconsole"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/remote"
	"github.com/ob-vss-ss19/blatt-3-sudo/messages"
	"github.com/ob-vss-ss19/blatt-3-sudo/tree"
	"log"
	"math/rand"
	"time"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

func RandStringRunes(n int) string {
	letterRunes := []rune(letters)

	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}

	return string(b)
}

type TreeData struct {
	pid    *actor.PID
	token  string
	marked bool
}

type RootActor struct {
	idToData map[int32]TreeData
	trees    []int32
	lastIns  int32
	behavior actor.Behavior
}

func newRoot() *RootActor {
	a := &RootActor{
		idToData: make(map[int32]TreeData),
		trees:    make([]int32, 0),
		behavior: actor.NewBehavior(),
	}
	a.behavior.Become(a.rootBehavior)

	return a
}

func (root *RootActor) Receive(ctx actor.Context) {
	root.behavior.Receive(ctx)
}

func (root *RootActor) rootBehavior(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *messages.CreateTreeRequest:
		log.Printf("Create Tree Request incoming! %v\n", msg)
		values := make(tree.Storage, msg.Capacity+1)
		pid := ctx.Spawn(actor.PropsFromProducer(func() actor.Actor {
			return tree.NewNode(int(msg.Capacity), &values)
		}))
		token := RandStringRunes(5)

		id := root.lastIns + 1
		root.trees = append(root.trees, id)
		root.lastIns++

		root.idToData[id] = TreeData{token: token, pid: pid, marked: false}

		ctx.Respond(&messages.CreateTreeResponse{
			TreeId: &messages.TreeIdentifier{
				Id:    id,
				Token: token,
			},
		})
	case *messages.ListTreesRequest:
		log.Printf("List Trees Request incoming! %v\n", msg)
		var results = make([]*messages.TreeIdentifier, 0, len(root.trees))

		for _, id := range root.trees {
			if id != 0 {
				if _, ok := root.idToData[id]; ok {
					results = append(results, &messages.TreeIdentifier{
						Id:    id,
						Token: "",
					})
				}
			}
		}

		ctx.Respond(&messages.ListTreesResponse{
			TreeIds: results,
		})
	case *messages.DeleteTreeRequest:
		log.Printf("Tree Deletion Request incoming! %v\n!", msg)
		if root.idToData[msg.TreeId.Id].token == msg.TreeId.Token {
			if root.idToData[msg.TreeId.Id].marked {
				ctx.Respond(&messages.DeleteTreeResponse{
					Success:           true,
					MarkedForDeletion: false,
				})
				ctx.Forward(root.idToData[msg.TreeId.Id].pid)
				root.idToData[msg.TreeId.Id].pid.Poison()
				delete(root.idToData, msg.TreeId.Id)

			} else {
				ctx.Respond(&messages.DeleteTreeResponse{
					Success:           true,
					MarkedForDeletion: true,
				})

				data := root.idToData[msg.TreeId.Id]
				data.marked = true
				root.idToData[msg.TreeId.Id] = data
			}
		} else {
			log.Printf("Token is wrong!\n")
		}
	case *messages.InsertRequest:
		log.Printf("Insert Request incoming %v ", msg)
		err := forward(ctx, root, msg.TreeId)
		if err != nil {
			ctx.Respond(&messages.InsertResponse{
				Success:      false,
				ErrorMessage: err.Error(),
			})
		}
	case *messages.SearchRequest:
		log.Printf("Search Request incoming %v ", msg)
		err := forward(ctx, root, msg.TreeId)
		if err != nil {
			ctx.Respond(&messages.SearchResponse{
				Success:      false,
				Entry:        nil,
				ErrorMessage: err.Error(),
			})
		}
	case *messages.RemoveRequest:
		log.Printf("Remove Request incoming %v ", msg)
		err := forward(ctx, root, msg.TreeId)
		if err != nil {
			ctx.Respond(&messages.RemoveResponse{
				Success:      false,
				RemovedPair:  nil,
				ErrorMessage: err.Error(),
			})
		}
	case *messages.TraverseRequest:
		log.Printf("Traverse Request incoming %v ", msg)
		err := forward(ctx, root, msg.TreeId)
		if err != nil {
			ctx.Respond(&messages.TraverseResponse{
				Success:      false,
				Pairs:        nil,
				ErrorMessage: err.Error(),
			})
		}
	}
}

func forward(ctx actor.Context, root *RootActor, data *messages.TreeIdentifier) error {
	treeData, ok := root.idToData[data.Id]

	if !ok {
		return fmt.Errorf("unknown tree identifier %d", data.Id)
	}

	if data.Token != treeData.token {
		return errors.New("the tree access token you supplied is incorrect")
	}

	ctx.Forward(treeData.pid)
	log.Printf("\n")

	return nil
}

func main() {
	printHeader()

	bind := flag.String(
		"bind",
		fmt.Sprintf("%s:%d", defaultName, defaultPort),
		fmt.Sprintf("to what name and address to bind the service: for example --bind=\"%s:%d\"", defaultName, defaultPort),
	)

	actorName := flag.String(
		"actor-name",
		defaultActorName,
		"the name of the tree service actor",
	)

	flag.Parse()

	remote.Start(*bind)
	rand.Seed(time.Now().UnixNano())

	var rootContext = actor.EmptyRootContext
	props := actor.PropsFromProducer(func() actor.Actor {
		return newRoot()
	})

	serverPID, err := rootContext.SpawnNamed(props, *actorName)
	if err != nil {
		log.Fatalf("Could not create root actor")
	}

	_, _ = console.ReadLine() // Wait for console input to terminate the application
	serverPID.GracefulPoison()
}

func printHeader() {
	fmt.Printf("%s\n\n", welcomeHeader)
}
