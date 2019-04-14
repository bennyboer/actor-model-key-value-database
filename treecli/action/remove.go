package action

import (
	"errors"
	"fmt"
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

func (Remove) Execute(client messages.TreeServiceClient, flags *util.Flags, args []string) error {
	log.Println("EXECUTE: Remove key-value pair")

	if flags.Id < 0 {
		return errors.New("please supply a valid tree ID")
	}
	if len(flags.Token) == 0 {
		return errors.New("please supply a valid Token")
	}
	if len(args) != 1 {
		return errors.New("the remove action expects only a key in the form: remove [key]")
	}

	// Parse key
	key, e := strconv.ParseInt(args[0], 10, 64)
	if e != nil {
		return errors.New(fmt.Sprintf("the key %s is not an integer", args[0]))
	}

	log.Printf("Key: %d\n", key)

	// TODO

	return nil
}
