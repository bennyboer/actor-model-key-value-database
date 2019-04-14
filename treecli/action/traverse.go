package action

import (
	"errors"
	"github.com/ob-vss-ss19/blatt-3-sudo/messages"
	"github.com/ob-vss-ss19/blatt-3-sudo/treecli/util"
	"log"
)

// Action traversing a tree.
type Traverse struct{}

func (Traverse) Identifier() string {
	return traverse
}

func (Traverse) Execute(client messages.TreeServiceClient, flags *util.Flags, args []string) error {
	log.Println("EXECUTE: Traverse tree")

	if flags.Id < 0 {
		return errors.New("please supply a valid tree ID")
	}
	if len(flags.Token) == 0 {
		return errors.New("please supply a valid Token")
	}

	// TODO

	return nil
}
