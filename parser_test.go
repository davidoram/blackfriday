package blackfriday

import (
	"testing"
)

func TestParse(t *testing.T) {
	input := "flower: 192.168.255.199 is HOST1"
	params := map[string]string{
		"ip":   "192.168.255.199",
		"host": "HOST1",
	}
	pc := Parse(input)
	if pc == nil {
		t.Errorf("Failed to parse '%s'", input)
	} else {
		assertParsedCommand(t, pc, "host_alias", params, input)
	}

}

func assertParsedCommand(t *testing.T, pc *ParsedCommand, command string, params map[string]string, input string) {

	if pc.command != command {
		t.Errorf("Command mismatch(%s). Expected %s, got %s", input, command, pc.command)
	}
	for key, value := range params {
		if pc.params[key] != value {
			t.Errorf("Paramater mismatch(%s). Expected %s, got %s for key %s", input, pc.params[key], value, key)
		}
	}
}
