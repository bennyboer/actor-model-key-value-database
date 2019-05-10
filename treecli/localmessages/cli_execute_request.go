package localmessages

import (
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/ob-vss-ss19/blatt-3-sudo/treecli/util"
	"strings"
)

// Execute request for the CLI actor.
type CLIExecuteRequest struct {
	// Arguments to execute.
	Arguments []string

	// Flags to augment the execution behavior.
	Flags *util.Flags

	/// PID of the remote actor.
	RemotePID *actor.PID
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
