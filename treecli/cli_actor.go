package main

import (
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/ob-vss-ss19/blatt-3-sudo/messages"
	"github.com/ob-vss-ss19/blatt-3-sudo/treecli/action"
	"github.com/ob-vss-ss19/blatt-3-sudo/treecli/local_messages"
	"github.com/ob-vss-ss19/blatt-3-sudo/treecli/util"
	"log"
)

// Actor realizing the CLI.
type CLIActor struct {
	// Current behavior of the actor.
	behavior actor.Behavior

	// Sender to which to respond after the result came from the server.
	sender *actor.PID
}

// Create a new CLI actor.
func NewCLIActor() *CLIActor {
	a := CLIActor{
		behavior: actor.NewBehavior(),
	}

	a.behavior.Become(a.ExecuteState)

	return &a
}

func (a *CLIActor) Receive(ctx actor.Context) {
	a.behavior.Receive(ctx)
}

// State in which the CLI actor is awaiting a command execution request.
func (a *CLIActor) ExecuteState(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *local_messages.CLIExecuteRequest:
		log.Printf("Received CLI Execution Request: %v\n", msg)

		a.sender = ctx.Sender()
		a.executeCommand(ctx, msg.Arguments, msg.Flags)
		a.behavior.Become(a.ReplyState)
	}
}

// State in which the CLI actor is awaiting an answer from the tree service.
func (a *CLIActor) ReplyState(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *messages.ListTreesResponse:
		ctx.Send(a.sender, local_messages.CLIExecuteReply{Message: fmt.Sprintf("%v", msg.TreeIds)})
	case *messages.CreateTreeResponse:
		ctx.Send(a.sender, local_messages.CLIExecuteReply{Message: fmt.Sprintf("%v", msg.TreeId)})
	case *messages.DeleteTreeResponse:
		var message string

		if msg.Success {
			if msg.MarkedForDeletion {
				message = "Marked tree for deletion. Execute command again to delete tree forever."
			} else {
				message = "Tree has been deleted."
			}
		} else {
			message = "Tree could not be deleted."
		}

		ctx.Send(a.sender, local_messages.CLIExecuteReply{Message: message})
	case *messages.InsertResponse:
		ctx.Send(a.sender, local_messages.CLIExecuteReply{Message: fmt.Sprintf("%v", msg.Success)})
	case *messages.SearchResponse:
		ctx.Send(a.sender, local_messages.CLIExecuteReply{Message: fmt.Sprintf("%v", msg.Success)})
	case *messages.RemoveResponse:
		ctx.Send(a.sender, local_messages.CLIExecuteReply{Message: fmt.Sprintf("%v", msg.Success)})
	case *messages.TraverseResponse:
		ctx.Send(a.sender, local_messages.CLIExecuteReply{Message: fmt.Sprintf("%v", msg.Pairs)})
	}
}

// Execute the passed arguments and flags as command for the CLI.
func (a *CLIActor) executeCommand(ctx actor.Context, arguments []string, flags *util.Flags) {
	// Create mapping of argument name to action.
	actionMap := make(map[string]action.Action)
	for _, a := range action.Actions {
		actionMap[a.Identifier()] = a
	}

	// Find the correct action to execute.
	startArgument := arguments[0]
	if startAction, ok := actionMap[startArgument]; ok {
		remote := actor.NewPID(fmt.Sprintf("%s:%d", flags.RemoteName, flags.RemotePort), flags.RemoteActorName)

		// Found action -> Execute with remaining arguments and flags
		if e := startAction.Execute(ctx, flags, arguments[1:], remote); e == nil {
			log.Println("Successfully sent command to tree service")
		} else {
			log.Fatalf("An error occurred :(\n%s", e.Error())
		}
	} else {
		log.Fatalf("Could not understand the argument \"%s\"", startArgument)
	}
}
