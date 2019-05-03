package main

import (
	"flag"
	"fmt"
	"github.com/AsynkronIT/goconsole"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/remote"
	"github.com/ob-vss-ss19/blatt-3-sudo/messages"
	"log"
)

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

	markedForDeletionTest := false // TODO Remove

	var rootContext = actor.EmptyRootContext
	props := actor.PropsFromFunc(func(ctx actor.Context) {
		switch msg := ctx.Message().(type) {
		case *messages.ListTreesRequest:
			log.Printf("List Trees Request incoming! %v\n", msg)

			testResult := []*messages.TreeIdentifier{
				{
					Id:    123,
					Token: "abc123f",
				},
				{
					Id:    845,
					Token: "73fbw93",
				},
			}[:]

			ctx.Respond(&messages.ListTreesResponse{
				TreeIds: testResult,
			})
		case *messages.DeleteTreeRequest:
			log.Printf("Tree deletion!")

			if markedForDeletionTest {
				ctx.Respond(&messages.DeleteTreeResponse{
					Success:           true,
					MarkedForDeletion: false,
				})

				markedForDeletionTest = false
			} else {
				ctx.Respond(&messages.DeleteTreeResponse{
					Success:           true,
					MarkedForDeletion: true,
				})

				markedForDeletionTest = true
			}
		}
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
