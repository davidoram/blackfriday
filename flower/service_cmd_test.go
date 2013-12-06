package flower

import (
	"testing"
)

func Test_HostIsLocalhost_True(t *testing.T) {
	hostname := "localhost"
	match, err := HostIsLocalhost(hostname)
	if err != nil {
		t.Errorf("HostIsLocalhost('%s') returned error: %v", hostname, err)
	}
	if !match {
		t.Errorf("HostIsLocalhost('%s') should return false", hostname)
	}
}

func Test_HostIsLocalhost_False(t *testing.T) {
	hostname := "www.gogle.com"
	match, err := HostIsLocalhost(hostname)
	if err != nil {
		t.Errorf("HostIsLocalhost('%s') returned error: %v", hostname, err)
	}
	if match {
		t.Errorf("HostIsLocalhost('%s') should return false", hostname)
	}
}

func Test_HostIsLocalhost_Error(t *testing.T) {
	hostname := "there_is_no_such_hast_as_this"
	_, err := HostIsLocalhost(hostname)
	if err == nil {
		t.Errorf("HostIsLocalhost('%s') should returned an error", hostname)
	}
}
