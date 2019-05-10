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

// Action removing a key-value pair from a tree.
type Remove struct{}

func (Remove) Identifier() string {
	return remove
}

func (Remove) Execute(ctx actor.Context, flags *util.Flags, args []string, remote *actor.PID) error {
	log.Println("EXECUTE: Remove key-value pair")

	if flags.ID < 0 {
		return errors.New("please supply a valid tree ID")
	}
	if len(flags.Token) == 0 {
		return errors.New("please supply a valid Token")
	}
	if len(args) != 1 {
		return errors.New("the remove action expects only a key in the form: remove [key]")
	}

	// Parse key to remove
	key, e := strconv.ParseInt(args[0], 10, 32)
	if e != nil {
		return errors.New(fmt.Sprintf("the key %s is not an integer", args[0]))
	}

	log.Printf("Remove key: %d\n", key)

	ctx.Request(remote, &messages.RemoveRequest{
		TreeId: &messages.TreeIdentifier{
			Id:    int32(flags.ID),
			Token: flags.Token,
		},
		Key: int32(key),
	})

	return nil
}
