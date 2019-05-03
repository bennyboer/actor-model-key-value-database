package action

import (
	"errors"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/ob-vss-ss19/blatt-3-sudo/messages"
	"github.com/ob-vss-ss19/blatt-3-sudo/treecli/util"
	"log"
)

/// Action to delete a tree.
type DeleteTree struct{}

func (DeleteTree) Identifier() string {
	return deleteTree
}

func (DeleteTree) Execute(ctx actor.Context, flags *util.Flags, args []string, remote *actor.PID) error {
	log.Println("EXECUTE: Delete tree")

	if flags.Id < 0 {
		return errors.New("please supply a valid tree ID")
	}
	if len(flags.Token) == 0 {
		return errors.New("please supply a valid Token")
	}

	ctx.Request(remote, &messages.DeleteTreeRequest{
		TreeId: &messages.TreeIdentifier{
			Id:    int32(flags.Id),
			Token: flags.Token,
		},
	})

	return nil
}
