package blackfriday

import (
	"bufio"
	"bytes"
	"container/list"
	"fmt"
	"net"
	"os"
	"regexp"
	"strconv"

//	"strings"
)

type Parser struct {
	hostname string
	commands *list.List
	errors   *list.List
	log      *os.File
	aliases  *list.List
}

// The FlowerParser evaluates code blocks for flower directives,
func NewParser() *Parser {
	parser := Parser{
		hostname: "Unknown",
		log:      os.Stderr,
		commands: list.New(),
		errors:   list.New(),
		// List of aliases this hostname is also known as
		aliases: list.New(),
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
			parser.commands.PushBack(command)
		}
	}
	return buffer.Bytes()
}

// Return markup for summary
func (parser *Parser) SummaryMarkup() []byte {
	buffer := new(bytes.Buffer)
	parser.ProcessAllCommands()
	buffer.WriteString("Hostname: " + parser.hostname)
	for e := parser.errors.Front(); e != nil; e = e.Next() {
		buffer.WriteString("\n\r")
		buffer.WriteString("ERROR: ")
		buffer.WriteString(e.Value.(error).Error())
	}
	return buffer.Bytes()
}

func (parser *Parser) ProcessHostAliasCommands() {
	for c := parser.commands.Front(); c != nil && c.Value.(*ParsedCommand).command == "host_alias"; c = c.Next() {
		pc := c.Value.(*ParsedCommand)
		fmt.Printf("Checking command :%+v\n", pc)
		intf, err := net.InterfaceAddrs()

		if err != nil {
			fmt.Println(err)
		} else {
			for _, ip := range intf {
				ipnet, ok := ip.(*net.IPNet)
				if ok {
					//fmt.Println("Checking IP :", ipnet.IP.String())
					if pc.params["ip"] == ipnet.IP.String() {
						fmt.Printf("\tLocal host matches IP %s, now aliased to %s\n", ipnet.IP.String(), pc.params["host"])
						parser.aliases.PushBack(pc.params["host"])
					}
				}
			}
		}
	}
}

// Process all of the commands we have gathered
func (parser *Parser) ProcessAllCommands() {
	parser.ProcessHostAliasCommands()
	parser.ProcessHostServiceCommands()
}

func (parser *Parser) CurrentHostIs(host string) bool {
	for e := parser.aliases.Front(); e != nil; e = e.Next() {
		if host == *(e.Value.(*string)) {
			return true
		}
	}
	return false
}

func (parser *Parser) ProcessHostServiceCommands() {
	// Process all host_service commands
	for c := parser.commands.Front(); c != nil; c = c.Next() {
		pc := c.Value.(*ParsedCommand)
		if pc.command == "host_service" {
			if parser.CurrentHostIs(pc.params["host"]) {
				fmt.Println("Check command :", pc)
			}
		}
	}
}

// The result of a parsed command
type ParsedCommand struct {
	command   string
	params    map[string]string
	satisfied bool
}

// Map from command name to a regular expression, that contains named
// subexpressions (?P<name>regex) that represent the arguments to the command
// that must be captured
var CommandRegex = map[string]*regexp.Regexp{
	// <ip> is <alias
	"host_alias": regexp.MustCompile(`^flower:\s*(?P<ip>[0-9\.\*]+)\s*is\s*(?P<host>\w+)\s*$`),
	// <alias> offers <service>
	"host_service": regexp.MustCompile(`^flower:\s*(?P<host>\w+)\s*offers\s*(?P<service>\w+)\s*$`),
}

// Map of host_alias -> ip address wildcard
var HostIp = map[string]string{}

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
			return BuildCommand(command, params)
		}
	}
	return nil
}

// Mappings from Services -> default port
var ServicePort = map[string]string{
	"http": "80",
}

// Build a ParsedCommand
func BuildCommand(command string, params map[string]string) *ParsedCommand {
	switch {
	case command == "host_service":
		if params["port"] == "" {
			params["port"] = ServicePort[params["service"]]
		}
		// check port is a int
		var _, err = strconv.Atoi(params["port"])
		if err == nil {
			return &ParsedCommand{
				command: command,
				params:  params,
			}
		}
		return nil
	}
	return &ParsedCommand{
		command: command,
		params:  params,
	}

}

// Evaluate a command
func RunCommand(pc *ParsedCommand) {
	switch {
	case pc.command == "host_alias":
		HostIp[pc.params["host"]] = pc.params["ip"]
		pc.satisfied = true
	case pc.command == "host_service":
		if RunningOn(pc.params["alias"]) {
			port, _ := strconv.Atoi(pc.params["port"])
			if PortOpen("127.0.0.1", port) {
				pc.satisfied = true
			} else {
				pc.satisfied = false
			}
		}
	}
}

// Determine if the program is running on a given host (alias)
func RunningOn(alias string) bool {
	return true
}

// Determine if a port is open on
func PortOpen(ip string, port int) bool {
	return true
}
