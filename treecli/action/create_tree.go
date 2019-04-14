package action

import (
	"github.com/ob-vss-ss19/blatt-3-sudo/messages"
	"github.com/ob-vss-ss19/blatt-3-sudo/treecli/util"
	"log"
)

/// Action to create a new, empty tree.
type CreateTree struct{}

func (CreateTree) Identifier() string {
	return createTree
}

func (CreateTree) Execute(client messages.TreeServiceClient, flags *util.Flags, args []string) error {
	log.Println("EXECUTE: Create tree")

	// TODO

	return nil
}
