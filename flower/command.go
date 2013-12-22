package flower

import ()

type Command interface {
	// Execute exectutes the command.
	Execute()

	// Return HTML id
	HtmlId() string

	// Return a string description of the command
	String() string
}
