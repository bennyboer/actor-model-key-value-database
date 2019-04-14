package main

import (
	"flag"
	"fmt"
	"github.com/ob-vss-ss19/blatt-3-sudo/messages"
	"github.com/ob-vss-ss19/blatt-3-sudo/treecli/action"
	"github.com/ob-vss-ss19/blatt-3-sudo/treecli/util"
	"google.golang.org/grpc"
	"log"
)

func main() {
	arguments := util.GetProgramArguments()
	flags := util.GetProgramFlags()

	log.Printf("Args: %v\n", arguments)
	log.Printf("Flags: %v\n", flags)

	// Set up service connection
	connection, err := grpc.Dial(fmt.Sprintf("%s:%d", flags.RemoteName, flags.RemotePort), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect to service: %v", err)
	}
	defer connection.Close()

	client := messages.NewTreeServiceClient(connection)

	if len(arguments) == 0 {
		printHelp()
	} else {
		// Create mapping of argument name to action.
		actionMap := make(map[string]action.Action)
		for _, a := range action.Actions {
			actionMap[a.Identifier()] = a
		}

		// Find the correct action to execute.
		startArgument := arguments[0]
		if startAction, ok := actionMap[startArgument]; ok {
			// Found action -> Execute with remaining arguments and flags
			if e := startAction.Execute(client, flags, arguments[1:]); e == nil {
				log.Println("Success")
			} else {
				log.Fatalf("An error occurred :(\n%s", e.Error())
			}
		} else {
			log.Printf("Could not understand the argument \"%s\"", startArgument)
			printHelp()
		}
	}
}

func printHelp() {
	fmt.Println("Actions:")
	for _, a := range action.Actions {
		fmt.Printf("- %s\n", a.Identifier())
	}

	fmt.Println("Flags:")
	flag.PrintDefaults()
}
