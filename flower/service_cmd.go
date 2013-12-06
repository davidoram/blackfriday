package flower

import (
//	"fmt"
	"github.com/anvie/port-scanner"
	"net"
//	"strings"
	"strconv"
)

// Constants used to indicate if a command applies to this host
const (
	NO_MATCH     = 0
	MATCH_HOST   = 1
	MATCH_CALLER = 2
)

// Constants used to indicate command success
const (
	ERROR = 0
	OK    = 1
	FAIL  = 2
)

type ServiceCommand struct {

	// Does this command match this host. See constants
	match int

	// Sucess flag. See constants
	result int

	// The services host, service name and port
	host    string
	service string
	port    int

	// The host that calls the service
	caller string

	// If an error occurs store it here
	err error
}

func (cmd *ServiceCommand) Execute() {

	var matchHost bool
	matchHost, cmd.err = HostIsLocalhost(cmd.host)
	if cmd.err != nil {
		return
	}

	var matchCaller bool
	matchCaller, cmd.err = HostIsLocalhost(cmd.caller)
	if cmd.err != nil {
		return
	}

	var ps *portscanner.PortScanner

	if matchHost {
		cmd.match = MATCH_HOST
		ps = portscanner.NewPortScanner(cmd.host)
	} else if matchCaller {
		cmd.match = MATCH_CALLER
		ps = portscanner.NewPortScanner(cmd.caller)
	} else {
		cmd.match = NO_MATCH
		return
	}

	// Scanning
	openedPorts := ps.GetOpenedPort(cmd.port, cmd.port)
	if len(openedPorts) == 1 {
		cmd.result = OK
	} else {
		cmd.result = FAIL
	}
}

func (cmd *ServiceCommand) String() string {
	str := cmd.host + "/" + cmd.service + "/" + strconv.Itoa(cmd.port) + " "
	switch cmd.match {
		case MATCH_HOST: 		str = str + "MATCH_HOST"
		case MATCH_CALLER:  	str = str + "MATCH_CALLER"
		case NO_MATCH:			str = str + "NO_MATCH"
	}
	
	return str;
}

// Return true if host name passed in, is the local host machine, otherwise returns false
func HostIsLocalhost(host string) (bool, error) {
	addr, err := net.InterfaceAddrs()
	if err != nil {
		return false, err
	}

	host_ip, err := net.LookupIP(host)
	if err != nil {
		return false, err
	}
	for _, address := range addr {
		localhost_ip, _, err := net.ParseCIDR(address.String())
		if err != nil {
			return false, err
		}
		for _, ip := range host_ip {
			if ip.Equal(localhost_ip) {
				return true, nil
			}
		}
	}
	return false, nil
}