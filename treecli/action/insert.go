package action

import (
	"errors"
	"fmt"
	"github.com/ob-vss-ss19/blatt-3-sudo/treecli/util"
	"log"
	"strconv"
	"strings"
)

// Action inserting a key-value pair in a tree.
type Insert struct{}

func (Insert) Identifier() string {
	return insert
}

func (Insert) Execute(args []string, flags *util.Flags) error {
	log.Println("EXECUTE: Insert key-value pair")

	if flags.Id < 0 {
		return errors.New("please supply a valid tree ID")
	}
	if len(flags.Token) == 0 {
		return errors.New("please supply a valid Token")
	}
	if len(args) < 2 {
		return errors.New("the insert action expects a key and a value in the form: insert [key] [value]")
	}

	// Parse key
	key, e := strconv.ParseInt(args[0], 10, 64)
	if e != nil {
		return errors.New(fmt.Sprintf("the key %s is not an integer", args[0]))
	}

	// Collect value
	var sb strings.Builder
	for i, part := range args[1:] {
		if part[0] == '"' || part[0] == '\'' {
			part = part[1:]
		}
		if part[len(part)-1] == '"' || part[len(part)-1] == '\'' {
			part = part[:len(part)-1]
		}

		if i > 0 {
			sb.WriteRune(' ')
		}
		sb.WriteString(part)
	}
	value := sb.String()

	log.Printf("Key: %d, Value: %s\n", key, value)

	// TODO

	return nil
}
