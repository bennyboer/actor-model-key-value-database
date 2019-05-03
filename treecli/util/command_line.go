package util

import (
	"os"
)

// Rune marking a flag.
const flagRune = '-'

// Fetch all program arguments received via command line.
func GetProgramArguments() []string {
	arguments := make([]string, 0)
	parts := os.Args[1:]
	for i := 0; i < len(parts); i++ {
		argument := parts[i]

		if isArgument(&argument) {
			arguments = append(arguments, argument)
		}
	}

	return arguments
}

// Check if the passed argument is a argument.
// An argument is NOT a flag which would start with a minus symbol '-'.
func isArgument(argument *string) bool {
	arg := *argument
	return arg[0] != flagRune
}
