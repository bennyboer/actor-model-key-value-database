package action

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/ob-vss-ss19/blatt-3-sudo/treecli/util"
)

// Action executable with the CLI.
type Action interface {
	/// Get the identifier of the action.
	/// The action is callable using the CLI and the identifier.
	/// Example: Identifier = "action", which would be executable with: "tree-cli action".
	Identifier() string

	// Execute the action with the given flags and arguments.
	Execute(ctx actor.Context, flags *util.Flags, args []string, remote *actor.PID) error
}

// All available command line actions.
var Actions = []Action{
	&List{},
	&CreateTree{},
	&DeleteTree{},
	&Insert{},
	&Remove{},
	&Search{},
	&Traverse{},
}
