package blackfriday

import (
	"bufio"
	"bytes"
	"container/list"
	"errors"
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
	aliasToIPNet map[string] *net.IPNet
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
		aliasToIPNet : make( map[string] *net.IPNet),
	}
	// resolve hostname
	var err error
	parser.hostname, err = os.Hostname()
	if err != nil {
		parser.errors.PushBack(errors.New("Unable to resolve hostname"))
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
		command := parser.Parse(line)
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

// Process all of the commands we have gathered
func (parser *Parser) ProcessAllCommands() {
	parser.ProcessHostAliasCommands()
	parser.ProcessHostServiceCommands()
	parser.PrintSummary()
}

// Process all host_alias commands
func (parser *Parser) ProcessHostAliasCommands() {
	for c := parser.commands.Front(); c != nil; c = c.Next() {
		pc, ok := c.Value.(*HostAliasCommand)
		if ok {
			intf, err := net.InterfaceAddrs()

			if err != nil {
				fmt.Println(err)
			} else {
				for _, ip := range intf {
					ipnet, ok := ip.(*net.IPNet)
					if ok {
						if pc.ip == ipnet.IP.String() {
							pc.isLocalHost = true
							parser.aliases.PushBack(pc.host)
							parser.aliasToIPNet[pc.host] = ipnet
						}
					}
				}
			}
		}
	}
}

// Process all host_service commands
func (parser *Parser) ProcessHostServiceCommands() {
	for c := parser.commands.Front(); c != nil; c = c.Next() {
		pc, ok := c.Value.(*HostServiceCommand)
		if ok && parser.CurrentHostIs(pc.host) {
			pc.isLocalHost = true
			pc.available = parser.ServiceRunning(pc)
		}
	}
}

func (parser *Parser) TotalErrors() int {
	var errors int
	for c := parser.commands.Front(); c != nil; c = c.Next() {
		pc, ok := c.Value.(*HostServiceCommand)
		if ok && pc.isLocalHost && !pc.available {
			errors += 1
		}
	}
	errors += parser.errors.Len()
	return errors
}

// Print a summary of findings
func (parser *Parser) PrintSummary() {
	fmt.Println("=======================================")
	fmt.Println("Flower summary")
	fmt.Println("--------------\n")
	fmt.Println("Host\n----")
	fmt.Printf("hostname: %s\n", parser.hostname)
	fmt.Printf("aliases:  ")
	for e := parser.aliases.Front(); e != nil; e = e.Next() {
		if e != parser.aliases.Front() {
			fmt.Printf(", ")
		}
		fmt.Printf(e.Value.(string))
	}
	fmt.Printf("\n\nAlias to IPAddr\n--------------\n")
	fmt.Printf("%-10s : %-20s\n", "Alias", "IP Address")
	for k,v := range(parser.aliasToIPNet) {
		fmt.Printf("%-10s : %-20s\n", k, v.IP.String())
	}

	fmt.Printf("\n\nLocal Services\n--------------\n")
	fmt.Printf("%-10s : %-5s : %5s\n", "Name", "Port", "Available")
	for c := parser.commands.Front(); c != nil; c = c.Next() {
		pc, ok := c.Value.(*HostServiceCommand)
		if ok && pc.isLocalHost {
			fmt.Printf("%-10s : %-5d : %v\n", pc.service, pc.port, pc.available)
		}
	}

	if parser.errors.Len() > 0 {
		fmt.Printf("\n\nErrors\n------\n")
		for e := parser.errors.Front(); e != nil; e = e.Next() {
			fmt.Printf("Err is :%v\n", e.Value)
		}
	}

	fmt.Printf("\n\nSummary\n-------\n")
	fmt.Printf("Total errors: %d\n", parser.TotalErrors())

}

func (parser *Parser) CurrentHostIs(host string) bool {
	for e := parser.aliases.Front(); e != nil; e = e.Next() {
		if host == e.Value.(string) {
			return true
		}
	}
	return false
}

// Determine if a port is open on
func (parser *Parser) ServiceRunning(c *HostServiceCommand) bool {
	hostPort := net.JoinHostPort(parser.aliasToIPNet[c.host].IP.String(), strconv.Itoa(c.port)) 
	fmt.Println("Calling Dial..."+ hostPort)
	var conn net.Conn
	conn, c.err = net.Dial("tcp", hostPort)
	if c.err != nil {
		parser.errors.PushBack(c.err)
		return false
	}
	conn.Close()
	return true
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
	"host_service": regexp.MustCompile(`^flower:\s*(?P<host>\w+)\s*offers\s*(?P<service>\w+)(:(?P<port>\d+))?\s*$`),
}

// Map of host_alias -> ip address wildcard
var HostIp = map[string]string{}

// Parse a string, find a matching command, or nil
func (parser *Parser) Parse(line string) FlowerCommand {
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
				parser.errors.PushBack(err)
			}
			return cmd

		}
	}
	return nil
}

// Mappings from Services -> default port
var ServicePort = map[string]int{
	"http": 80,
	"smtp": 25,
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
