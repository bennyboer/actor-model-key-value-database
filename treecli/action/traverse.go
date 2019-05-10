package action

import (
	"errors"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/ob-vss-ss19/blatt-3-sudo/messages"
	"github.com/ob-vss-ss19/blatt-3-sudo/treecli/util"
	"log"
)

// Action traversing a tree.
type Traverse struct{}

func (Traverse) Identifier() string {
	return traverse
}

func (Traverse) Execute(ctx actor.Context, flags *util.Flags, args []string, remote *actor.PID) error {
	log.Println("EXECUTE: Traverse tree")

	if flags.ID < 0 {
		return errors.New("please supply a valid tree ID")
	}
	if len(flags.Token) == 0 {
		return errors.New("please supply a valid Token")
	}

	ctx.Request(remote, &messages.TraverseRequest{
		TreeId: &messages.TreeIdentifier{
			Id:    int32(flags.ID),
			Token: flags.Token,
		},
	})

	return nil
}
