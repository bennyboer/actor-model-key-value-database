package action

import (
	"errors"
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/ob-vss-ss19/blatt-3-sudo/messages"
	"github.com/ob-vss-ss19/blatt-3-sudo/treecli/util"
	"log"
	"strconv"
)

/// Action to create a new, empty tree.
type CreateTree struct{}

func (CreateTree) Identifier() string {
	return createTree
}

func (CreateTree) Execute(ctx actor.Context, flags *util.Flags, args []string, remote *actor.PID) error {
	log.Println("EXECUTE: Create tree")

	// Parse capacity
	capacity, e := strconv.ParseInt(args[0], 10, 32)
	if e != nil {
		return errors.New(fmt.Sprintf("the capacity %s is not an integer", args[0]))
	}

	ctx.Request(remote, &messages.CreateTreeRequest{
		Capacity: int32(capacity),
	})

	return nil
}
