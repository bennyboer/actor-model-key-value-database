package local_messages

import (
	"fmt"
	"github.com/ob-vss-ss19/blatt-3-sudo/treecli/util"
	"strings"
)

// Execute request for the CLI actor.
type CLIExecuteRequest struct {
	// Arguments to execute.
	Arguments []string

	// Flags to augment the execution behavior.
	Flags *util.Flags
}

func (r *CLIExecuteRequest) String() string {
	var sb strings.Builder

	sb.WriteString("Arguments: ")
	sb.WriteString(fmt.Sprintf("%v", r.Arguments))
	sb.WriteString(", Flags: ")
	sb.WriteString(fmt.Sprintf("%v", r.Flags))
	sb.WriteRune('\n')

	return sb.String()
}
