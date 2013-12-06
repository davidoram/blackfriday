package flower

import ()

type Command interface {
	// Execute exectutes the command.
	Execute()

	// Return a string description of the command
	String() string
}
