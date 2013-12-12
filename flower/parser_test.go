package flower

import (
	"testing"
)

func Test_OffersService(t *testing.T) {
	input := "flower: HOST1 offers http"
	AssertEquals(ServiceCommand{host: "HOST1", service: "http", port: 80, caller: "HOST1"}, input, t)

	input = "flower: HOST1 offers http:8080"
	AssertEquals(ServiceCommand{host: "HOST1", service: "http", port: 8080, caller: "HOST1"}, input, t)

	input = "flower: my.host.name offers http:8080"
	AssertEquals(ServiceCommand{host: "my.host.name", service: "http", port: 8080, caller: "my.host.name"}, input, t)
}

func Test_UsesService(t *testing.T) {
	input := "flower: HOST1 uses http at HOST2"
	AssertEquals(ServiceCommand{host: "HOST2", service: "http", port: 80, caller: "HOST1"}, input, t)

	input = "flower: HOST1 uses http:8080 at HOST2"
	AssertEquals(ServiceCommand{host: "HOST2", service: "http", port: 8080, caller: "HOST1"}, input, t)

	input = "flower: my.host.name uses http:8080 at my.remote.host"
	AssertEquals(ServiceCommand{host: "my.remote.host", service: "http", port: 8080, caller: "my.host.name"}, input, t)
}

func AssertEquals(expected ServiceCommand, input string, t *testing.T) {
	actual := Parse(input)
	a := actual.(*ServiceCommand)
	if expected.host != a.host {
		t.Errorf("Host mismatch. Expected %s, got %s. Input %s", expected.host, a.host, input)
	}
	if expected.port != a.port {
		t.Errorf("Port mismatch. Expected %s, got %s. Input %s", expected.port, a.port, input)
	}
	if expected.service != a.service {
		t.Errorf("Service mismatch. Expected %s, got %s. Input %s", expected.service, a.service, input)
	}
	if expected.caller != a.caller {
		t.Errorf("Caller mismatch. Expected %s, got %s. Input %s", expected.caller, a.caller, input)
	}
}
