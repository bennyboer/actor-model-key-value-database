package action

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/ob-vss-ss19/blatt-3-sudo/messages"
	"github.com/ob-vss-ss19/blatt-3-sudo/treecli/util"
	"log"
)

// Action lists all tree IDs.
type List struct{}

func (List) Identifier() string {
	return list
}

func (List) Execute(ctx actor.Context, flags *util.Flags, args []string, remote *actor.PID) error {
	log.Println("EXECUTE: List trees")

	ctx.Request(remote, &messages.ListTreesRequest{})

	return nil
}
