package action

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/ob-vss-ss19/blatt-3-sudo/messages"
	"github.com/ob-vss-ss19/blatt-3-sudo/treecli/util"
	"log"
)

/// Action to create a new, empty tree.
type CreateTree struct{}

func (CreateTree) Identifier() string {
	return createTree
}

func (CreateTree) Execute(ctx actor.Context, flags *util.Flags, args []string, remote *actor.PID) error {
	log.Println("EXECUTE: Create tree")

	ctx.Request(remote, &messages.CreateTreeRequest{})

	return nil
}
