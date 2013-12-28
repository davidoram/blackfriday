package flower

import (
	"bytes"
	"container/list"
	"fmt"
//	"strings"
)

type Interpreter interface {
	// Evaluate markdown code
	EvaluateCode(markdown_code string)

	// Return a report
	SummaryReport() string
}

type StandardInterpreter struct {
	commands *list.List
}

// Create and return an Interpreter
func NewInterpreter() *StandardInterpreter {
	interpreter := StandardInterpreter{
		commands: list.New(),
	}
	return &interpreter
}

// Evaluate a chunk of Markdown code. If it contains flower directives then evaluate it
func (interpreter *StandardInterpreter) EvaluateCode(line string) Command {
	command := Parse(line)
	if command != nil {
		interpreter.commands.PushBack(command)
		command.Execute()
	}
	return command
}


// Return a summary of findings
func (interpreter *StandardInterpreter) SummaryReport() []byte {
	buf := bytes.NewBufferString("")
	fmt.Fprintln(buf, "<table><tr><th>Flower Summary</th></tr>")

	fmt.Fprintln(buf, "<tr><th>Rules</th></tr>")
	for element := interpreter.commands.Front(); element != nil; element = element.Next() {
		cmd, ok := element.Value.(*ServiceCommand)
		if ok {
			fmt.Fprintln(buf, "<tr><td>", cmd.String(), "</tr></td>")
		}
	}
	fmt.Fprintln(buf, "</table>")
	return buf.Bytes()
}
