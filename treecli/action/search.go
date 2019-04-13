package action

import (
	"errors"
	"fmt"
	"github.com/ob-vss-ss19/blatt-3-sudo/treecli/util"
	"log"
	"strconv"
)

// Action searching for key-value pair in a tree.
type Search struct{}

func (Search) Identifier() string {
	return search
}

func (Search) Execute(args []string, flags *util.Flags) error {
	log.Println("EXECUTE: Search key-value pair in tree")

	if flags.Id < 0 {
		return errors.New("please supply a valid tree ID")
	}
	if len(flags.Token) == 0 {
		return errors.New("please supply a valid Token")
	}
	if len(args) != 1 {
		return errors.New("the search action expects only a key in the form: search [key]")
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
