package blackfriday

import (
	"bufio"
	"bytes"
	"container/list"
	"fmt"
	"os"
)

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
// flower: app_server -> db_server:mysql
// When run on db_server, will check that a server is listening on the standard mysql port
// When run on app_server, will a remote connection can be established to db_server on the standard mysql port
//
// flower: host <hostname> == ip
// eg:
// flower: app_server == 127.0.0.1
// Will check that app_server has ip address 127.0.0.1
//
// User/Group based checks:
// flower: group <group> exists
// flower: user <user> exists
// flower: user <user> in <group>

func FlowerParser() *Parser {
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

	scanner := bufio.NewScanner(bytes.NewReader(text))
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(parser.log, "GOT LINE: "+line)
	}
	return text
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
