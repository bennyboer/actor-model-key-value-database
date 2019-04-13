package util

import (
	"flag"
	"fmt"
)

// Flags the CLI is able to understand.
type Flags struct {
	// Token of the tree service session
	Token string

	// Id of the tree to work with
	Id int
}

/// Get all program flags received via command line.
func GetProgramFlags() *Flags {
	token := flag.String("token", "", "Token of the tree service session")
	id := flag.Int("id", -1, "Id of the tree to work with")

	flag.Parse()

	return &Flags{
		Token: *token,
		Id:    *id,
	}
}

func (f *Flags) String() string {
	return fmt.Sprintf("{ token: %s, id: %v }", f.Token, f.Id)
}
