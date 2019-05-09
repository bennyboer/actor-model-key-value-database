package main

import (
	"flag"
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/remote"
	"github.com/ob-vss-ss19/blatt-3-sudo/treecli/action"
	"github.com/ob-vss-ss19/blatt-3-sudo/treecli/local_messages"
	"github.com/ob-vss-ss19/blatt-3-sudo/treecli/util"
	"log"
)

func main() {
	arguments := util.GetProgramArguments()
	flags := util.GetProgramFlags()

	if len(arguments) == 0 {
		printHelp()
	} else {
		result, err := process(arguments, flags).Result()
		if err != nil {
			log.Fatalf("Command execution failed:\n%s\n", err.Error())
		}

		response, ok := result.(*local_messages.CLIExecuteReply)
		if !ok {
			log.Fatalf("Answer of the CLI Actor is incorrect")
		}

		log.Println(response.Message)
	}
}

// Start the CLI actor and process the passed arguments and flags.
func process(arguments []string, flags *util.Flags) *actor.Future {
	remote.Start(fmt.Sprintf("%s:%d", flags.Name, flags.Port)) // Register as remote actor

	actorProps := actor.PropsFromProducer(func() actor.Actor {
		return NewCLIActor()
	})

	var rootContext = actor.EmptyRootContext
	cliActor := rootContext.Spawn(actorProps)

	remotePID := actor.NewPID(fmt.Sprintf("%s:%d", flags.RemoteName, flags.RemotePort), flags.RemoteActorName)

	future := rootContext.RequestFuture(
		cliActor,
		&local_messages.CLIExecuteRequest{
			Arguments: arguments,
			Flags:     flags,
			RemotePID: remotePID,
		},
		flags.Timeout,
	)

	return future
}

func printHelp() {
	fmt.Println("Actions:")
	for _, a := range action.Actions {
		fmt.Printf("- %s\n", a.Identifier())
	}

	fmt.Println("Flags:")
	flag.PrintDefaults()
}
