package action

import (
	"context"
	"github.com/ob-vss-ss19/blatt-3-sudo/messages"
	"github.com/ob-vss-ss19/blatt-3-sudo/treecli/util"
	"io"
	"log"
)

// Action lists all tree IDs.
type List struct{}

func (List) Identifier() string {
	return list
}

func (List) Execute(client messages.TreeServiceClient, flags *util.Flags, args []string) error {
	log.Println("EXECUTE: List trees")

	ctx, cancel := context.WithTimeout(context.Background(), util.DefaultTimeout)
	defer cancel()

	stream, e := client.ListTrees(ctx, &messages.VoidMessage{})
	if e != nil {
		return e
	}

	log.Println("Available Tree IDs: ---")
	for {
		treeId, err := stream.Recv()

		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		log.Printf("  - %d", treeId.Id)
	}
	log.Println("-----------------------")

	return nil
}
