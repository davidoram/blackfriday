package blackfriday

import (
	"bufio"
	"bytes"
	"container/list"
	"os"
	"regexp"
)

const ParsedOkSymbol = " - OK"

type Parser struct {
	hostname string
	errors   *list.List
	log      *os.File
}

// The FlowerParser evaluates code blocks for flower directives,
// executes them and returns the evaluated data as text
// It understands the following syntax:
//
// flower: host <hostname|ip> -> <hostname:ip>:<port number|port name>
// eg:
// flower: host app_server -> db_server:mysql
// When run on db_server, will check that a server is listening on the standard mysql port
// When run on app_server, will a remote connection can be established to db_server on the standard mysql port
//
// flower: host <hostname> == ip
// eg:
// flower: host app_server == 127.0.0.1
// Will check that app_server has ip address 127.0.0.1
//
// User/Group based checks:
// flower: group <group> exists
// flower: user <user> exists
// flower: user <user> in <group>

func NewParser() *Parser {
	parser := Parser{
		hostname: "Unknown",
		log:      os.Stderr,
		errors:   list.New(),
	}
	// resolve hostname
	var err error
	parser.hostname, err = os.Hostname()
	if err != nil {
		parser.errors.PushBack("Unable to resolve hostname")
	}

	return &parser
}

// Evaluate text, code and return text altered
func (parser *Parser) EvaluateCode(text []byte) []byte {

	buffer := new(bytes.Buffer)
	scanner := bufio.NewScanner(bytes.NewReader(text))
	for scanner.Scan() {
		line := scanner.Text()
		buffer.WriteString(line)
		command := Parse(line)
		if command != nil {
			//buffer.WriteString("<br>")
			buffer.WriteString(Evaluate(command))
			//buffer.WriteString("<br>")
		}
	}
	return buffer.Bytes()
}

// Return markup for summary
func (parser *Parser) SummaryMarkup() []byte {
	buffer := new(bytes.Buffer)
	buffer.WriteString("Hostname: " + parser.hostname)
	for e := parser.errors.Front(); e != nil; e = e.Next() {
		buffer.WriteString("\n\r")
		buffer.WriteString("ERROR: ")
		buffer.WriteString(e.Value.(error).Error())
	}
	return buffer.Bytes()
}

// The result of a parsed command
type ParsedCommand struct {
	command string
	params  map[string]string
}

// Map from command name to a regular expression, that contains named
// subexpressions (?P<name>regex) that represent the arguments to the command
// that must be captured
var CommandRegex = map[string]*regexp.Regexp{
	// <ip> is <alias
	"host_alias": regexp.MustCompile(`^flower:\s*(?P<ip>[0-9\.\*]+)\s*is\s*(?P<host>\w+)\s*$`),
}

// Map of host_alias -> ip address wildcard
var HostIp = map[string]string {}

// Parse a string, find a matching command, or nil
func Parse(line string) *ParsedCommand {
	//fmt.Printf("Parse : %s length  %d\n",line, len(line))
	for command, regex := range CommandRegex {
		//fmt.Printf("Is a  %s ? %t\n",command, regex.MatchString(line))
		if regex.MatchString(line) {
			// Turn parenthesized sub expressions into a map[string]string
			params := make(map[string]string)
			for _, key := range regex.SubexpNames() {
				params[key] = regex.ReplaceAllString(line, "${"+key+"}")
			}
			//fmt.Printf("Got params %q\n", params)
			return &ParsedCommand{
				command: command,
				params:  params,
			}
		}
	}
	return nil
}

// Evaluate a command
func Evaluate(pc *ParsedCommand) string {
	switch {
    case pc.command == "host_alias":
        HostIp[pc.params["host"]] = pc.params["ip"]
        return ParsedOkSymbol
    }	
    return ""
}
