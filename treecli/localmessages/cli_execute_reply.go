package localmessages

// Reply from the CLI actor.
type CLIExecuteReply struct {
	// The result message
	Message string

	// The original result structure
	Original interface{}
}
