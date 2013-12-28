package flower

import ()

type Command interface {
	// Execute exectutes the command.
	Execute()

	// Return HTML class(es) to be applied to a command
	HtmlClass() string

	// Return a string description of the command
	String() string
}
