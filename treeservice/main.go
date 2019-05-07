package main

import (
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

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
var root *RootActor = nil

func RandStringRunes(n int) string {
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
	lastIns  int
}

func getRoot() *RootActor {
	if root == nil {
		root = &RootActor{
			idToData: make(map[int32]TreeData),
			trees:    make([]int32, 2),
		}
	}
	return root
}

func forward(ctx actor.Context, root *RootActor, data *messages.TreeIdentifier) {
	if root.idToData[data.Id].token == data.Token {
		ctx.Forward(root.idToData[data.Id].pid)
		log.Printf("\n")
	} else {
		log.Printf("Wrong token!\n")
	}
}

func rootBehavior(ctx actor.Context) {
	var root = getRoot()

	switch msg := ctx.Message().(type) {
	case *messages.CreateTreeRequest:
		log.Printf("Create Tree Request incoming! %v\n", msg)
		pid := ctx.Spawn(actor.PropsFromProducer(tree.NewNode))
		token := RandStringRunes(5)

		var id int32
		if root.lastIns == 0 {
			id = 1
		} else {
			id = root.trees[root.lastIns] + 1
			root.lastIns++
		}
		root.trees[root.lastIns] = id

		root.idToData[id] = TreeData{token: token, pid: pid, marked: false}

		ctx.Respond(&messages.CreateTreeResponse{
			TreeId: &messages.TreeIdentifier{
				Id:    id,
				Token: token,
			},
		})
	case *messages.ListTreesRequest:
		log.Printf("List Trees Request incoming! %v\n", msg)
		var results = make([]*messages.TreeIdentifier, len(root.trees))

		for _, id := range root.trees {
			results = append(results, &messages.TreeIdentifier{
				Id:    id,
				Token: root.idToData[id].token,
			})
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
		forward(ctx, root, msg.TreeId)

	case *messages.SearchRequest:
		log.Printf("Search Request incoming %v ", msg)
		forward(ctx, root, msg.TreeId)

	case *messages.RemoveRequest:
		log.Printf("Remove Request incoming %v ", msg)
		forward(ctx, root, msg.TreeId)

	case *messages.TraverseRequest:
		log.Printf("Traverse Request incoming %v ", msg)
		forward(ctx, root, msg.TreeId)
	}
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
	props := actor.PropsFromFunc(rootBehavior)
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
