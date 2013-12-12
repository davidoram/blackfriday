package flower

import (
	"fmt"
	"regexp"
	"strconv"
)

// Mappings from Services -> default port
var ServicePort = map[string]int{
	"http": 80,
	"smtp": 25,
}

// Map from command name to a regular expression, that contains named
// subexpressions (?P<name>regex) that represent the arguments to the command
// that must be captured
var CommandRegex = map[string]*regexp.Regexp{
	// <ip> is <alias
	//"host_alias": regexp.MustCompile(`^flower:\s*(?P<ip>[0-9\.\*]+)\s*is\s*(?P<host>\w+)\s*$`),
	// <alias> offers <service>:port
	"local_service": regexp.MustCompile(`^flower:\s*(?P<host>\S+)\s*offers\s*(?P<service>\w+)(:(?P<port>\d+))?\s*$`),
	// <alias> uses <service>:port at <alias>
	"remote_service": regexp.MustCompile(`^flower:\s*(?P<local>\S+)\s*uses\s*(?P<service>\w+)(:(?P<port>\d+))?\s+(at)?\s+(?P<remote>\S+)\s*$`),
}

// Parse a string, find a matching command, or nil
func Parse(line string) Command {
	fmt.Println("Parsing: " + line)
	for command, regex := range CommandRegex {
	fmt.Print("Matches " + command + " : ")
		//fmt.Printf("Is a  %s ? %t\n",command, regex.MatchString(line))
		if regex.MatchString(line) {
			fmt.Println("Yes")
			// Turn parenthesized sub expressions into a map[string]string
			params := make(map[string]string)
			for _, key := range regex.SubexpNames() {
				params[key] = regex.ReplaceAllString(line, "${"+key+"}")
			}
			//fmt.Printf("Got params %q\n", params)
			return BuildCommand(command, params)
		} else {
			fmt.Println("No")
		}
	}
	return nil
}

// Build a ParsedCommand
func BuildCommand(command string, params map[string]string) Command {
	switch {
	case command == "local_service":
		var port int
		var err error
		if params["port"] == "" {
			port = ServicePort[params["service"]]
		} else {
			// check port is a int
			port, err = strconv.Atoi(params["port"])
		}
		if err == nil {
			return &ServiceCommand{
				host:    params["host"],
				service: params["service"],
				port:    port,
				caller:  params["host"],
			}
		}
	case command == "remote_service":
		var port int
		var err error
		if params["port"] == "" {
			port = ServicePort[params["service"]]
		} else {
			// check port is a int
			port, err = strconv.Atoi(params["port"])
		}
		if err == nil {
			return &ServiceCommand{
				host:    params["remote"],
				service: params["service"],
				port:    port,
				caller:  params["local"],
			}
		}
	}
	return nil
}
