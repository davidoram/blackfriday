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
	for c := parser.commands.Front(); c != nil; c = c.Next() {
		pc, ok := c.Value.(*HostAliasCommand)
		if ok {
			fmt.Printf("Checking command :%+v\n", pc)
			intf, err := net.InterfaceAddrs()

			if err != nil {
				fmt.Println(err)
			} else {
				for _, ip := range intf {
					ipnet, ok := ip.(*net.IPNet)
					if ok {
						if pc.ip == ipnet.IP.String() {
							fmt.Printf("\tLocal host matches IP %s, now aliased to %s\n", ipnet.IP.String(), pc.host)
							pc.isLocalHost = true
							parser.aliases.PushBack(pc.host)
						}
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
		if host == e.Value.(string) {
			return true
		}
	}
	return false
}

func (parser *Parser) ProcessHostServiceCommands() {
	fmt.Println("\nProcessHostServiceCommands.")
	// Process all host_service commands
	for c := parser.commands.Front(); c != nil; c = c.Next() {
		pc, ok := c.Value.(*HostServiceCommand)

		if ok && parser.CurrentHostIs(pc.host) {
			pc.isLocalHost = true
			fmt.Println("Check command :", pc)
			pc.available = parser.ServiceRunning(pc)
			fmt.Printf("\tHost service is running: %+v\n", pc)
		}
	}
}

// Determine if a port is open on
func (parser *Parser) ServiceRunning(c *HostServiceCommand) bool {
	var conn net.Conn
	conn, c.err = net.Dial("tcp", c.host+":"+strconv.Itoa(c.port))
	if c.err == nil {
		conn.Close()
		return true
	}
	return false
}

// The result of a parsed command
type FlowerCommand interface {
}

type HostAliasCommand struct {
	ip          string
	host        string
	isLocalHost bool
}

type HostServiceCommand struct {
	host        string
	isLocalHost bool
	service     string
	port        int
	available   bool
	err         error
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
func Parse(line string) FlowerCommand {
	fmt.Printf("Parse : %s length  %d\n", line, len(line))
	for command, regex := range CommandRegex {
		//fmt.Printf("Is a  %s ? %t\n",command, regex.MatchString(line))
		if regex.MatchString(line) {
			// Turn parenthesized sub expressions into a map[string]string
			params := make(map[string]string)
			for _, key := range regex.SubexpNames() {
				params[key] = regex.ReplaceAllString(line, "${"+key+"}")
			}
			//fmt.Printf("Got params %q\n", params)
			cmd, err := BuildCommand(command, params)
			if err != nil {
				fmt.Println("ERROR: Parsing", err)
			} else {
				fmt.Printf("Parsing match %+v\n", cmd)
			}
			return cmd

		}
	}
	return nil
}

// Mappings from Services -> default port
var ServicePort = map[string]int{
	"http": 80,
}

// Build a ParsedCommand
func BuildCommand(command string, params map[string]string) (FlowerCommand, error) {
	switch {
	case command == "host_service":
		var port int
		var err error
		if params["port"] == "" {
			port = ServicePort[params["service"]]
		} else {
			// check port is a int
			port, err = strconv.Atoi(params["port"])
		}
		if err == nil {
			return &HostServiceCommand{
				host:    params["host"],
				service: params["service"],
				port:    port,
			}, nil
		}
	case command == "host_alias":
		return &HostAliasCommand{
			ip:   params["ip"],
			host: params["host"],
		}, nil
	}
	return nil, nil
}
