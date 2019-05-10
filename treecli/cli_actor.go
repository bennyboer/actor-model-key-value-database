package main

import (
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/ob-vss-ss19/blatt-3-sudo/messages"
	"github.com/ob-vss-ss19/blatt-3-sudo/treecli/action"
	"github.com/ob-vss-ss19/blatt-3-sudo/treecli/local_messages"
	"github.com/ob-vss-ss19/blatt-3-sudo/treecli/util"
	"log"
	"strings"
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
		a.executeCommand(ctx, msg.Arguments, msg.Flags, msg.RemotePID)
		a.behavior.Become(a.ReplyState)
	}
}

// State in which the CLI actor is awaiting an answer from the tree service.
func (a *CLIActor) ReplyState(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *messages.ListTreesResponse:
		var sb strings.Builder
		sb.WriteString("Available Tree IDs:\n")

		for _, idPtr := range msg.TreeIds {
			id := *idPtr

			sb.WriteString("  - ")
			sb.WriteString(fmt.Sprintf("%d\n", id.Id))
		}

		ctx.Send(a.sender, &local_messages.CLIExecuteReply{
			Message:  sb.String(),
			Original: msg,
		})

		a.behavior.Become(a.ExecuteState)
	case *messages.CreateTreeResponse:
		ctx.Send(a.sender, &local_messages.CLIExecuteReply{
			Message:  fmt.Sprintf("Created tree with ID: %d and Token: '%s'\n", msg.TreeId.Id, msg.TreeId.Token),
			Original: msg,
		})

		a.behavior.Become(a.ExecuteState)
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

		ctx.Send(a.sender, &local_messages.CLIExecuteReply{
			Message:  message,
			Original: msg,
		})

		a.behavior.Become(a.ExecuteState)
	case *messages.InsertResponse:
		var message string
		if msg.Success {
			message = "Inserted successfully."
		} else {
			message = "Insert did not work."
		}

		ctx.Send(a.sender, &local_messages.CLIExecuteReply{
			Message:  message,
			Original: msg,
		})

		a.behavior.Become(a.ExecuteState)
	case *messages.SearchResponse:
		var message string
		if msg.Success {
			message = fmt.Sprintf("Found key-value pair with Key: %d and Value '%s'.", msg.Entry.Key, msg.Entry.Value)
		} else {
			message = "Could not find key-value pair."
		}

		ctx.Send(a.sender, &local_messages.CLIExecuteReply{
			Message:  message,
			Original: msg,
		})

		a.behavior.Become(a.ExecuteState)
	case *messages.RemoveResponse:
		var message string
		if msg.Success {
			message = fmt.Sprintf("Removed key-value pair with Key: %d and Value '%s'.", msg.RemovedPair.Key, msg.RemovedPair.Value)
		} else {
			message = "Could not remove key-value pair."
		}

		ctx.Send(a.sender, &local_messages.CLIExecuteReply{
			Message:  message,
			Original: msg,
		})

		a.behavior.Become(a.ExecuteState)
	case *messages.TraverseResponse:
		var message string
		if msg.Pairs != nil {
			var sb strings.Builder

			sb.WriteString("All key-value pairs:\n")

			for _, pairPtr := range msg.Pairs {
				pair := *pairPtr
				sb.WriteString(fmt.Sprintf("  - Key: %d, Value: '%s'\n", pair.Key, pair.Value))
			}

			message = sb.String()
		} else {
			message = "Could not traverse tree."
		}

		ctx.Send(a.sender, &local_messages.CLIExecuteReply{
			Message:  message,
			Original: msg,
		})

		a.behavior.Become(a.ExecuteState)
	}
}

// Execute the passed arguments and flags as command for the CLI.
func (a *CLIActor) executeCommand(ctx actor.Context, arguments []string, flags *util.Flags, remotePID *actor.PID) {
	// Create mapping of argument name to action.
	actionMap := make(map[string]action.Action)
	for _, a := range action.Actions {
		actionMap[a.Identifier()] = a
	}

	// Find the correct action to execute.
	startArgument := arguments[0]
	if startAction, ok := actionMap[startArgument]; ok {
		// Found action -> Execute with remaining arguments and flags
		if e := startAction.Execute(ctx, flags, arguments[1:], remotePID); e == nil {
			log.Println("Successfully sent command to tree service")
		} else {
			log.Fatalf("An error occurred :(\n%s", e.Error())
		}
	} else {
		log.Fatalf("Could not understand the argument \"%s\"", startArgument)
	}
}
