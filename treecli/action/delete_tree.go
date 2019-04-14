package action

import (
	"errors"
	"github.com/ob-vss-ss19/blatt-3-sudo/messages"
	"github.com/ob-vss-ss19/blatt-3-sudo/treecli/util"
	"log"
)

/// Action to delete a tree.
type DeleteTree struct{}

func (DeleteTree) Identifier() string {
	return deleteTree
}

func (DeleteTree) Execute(client messages.TreeServiceClient, flags *util.Flags, args []string) error {
	log.Println("EXECUTE: Delete tree")

	if flags.Id < 0 {
		return errors.New("please supply a valid tree ID")
	}
	if len(flags.Token) == 0 {
		return errors.New("please supply a valid Token")
	}

	// TODO

	return nil
}
