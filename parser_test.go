package blackfriday

import (
	"testing"
)

func TestParseHostAlias(t *testing.T) {
	input := "flower: 192.168.255.199 is HOST1"
	AssertEqualsHostAlias(HostAliasCommand { host: "HOST1", ip: "192.168.255.199",}, input, t)
}

func TestParseService(t *testing.T) {
	input := "flower: HOST1 offers http"
	AssertEqualsHostService(HostServiceCommand { host: "HOST1", service: "http", port: 80,}, input, t)

	input = "flower: HOST1 offers http:8080"
	AssertEqualsHostService(HostServiceCommand { host: "HOST1", service: "http", port: 8080,}, input, t)
}

func AssertEqualsHostAlias(expected HostAliasCommand, input string, t *testing.T) {
	actual := Parse(input)
	a := actual.(*HostAliasCommand)
	if expected.host != a.host {
		t.Errorf("Host mismatch. Expected %s, got %s. Input %s", expected.host, a.host, input)
	}
	if expected.ip != a.ip {
		t.Errorf("IP mismatch. Expected %s, got %s. Input %s", expected.ip, a.ip, input)
	}
}

func AssertEqualsHostService(expected HostServiceCommand, input string, t *testing.T) {
	actual := Parse(input)
	a := actual.(*HostServiceCommand)
	if expected.host != a.host {
		t.Errorf("Host mismatch. Expected %s, got %s. Input %s", expected.host, a.host, input)
	}
	if expected.port != a.port {
		t.Errorf("Port mismatch. Expected %s, got %s. Input %s", expected.port, a.port, input)
	}
	if expected.service != a.service {
		t.Errorf("Service mismatch. Expected %s, got %s. Input %s", expected.service, a.service, input)
	}
}

