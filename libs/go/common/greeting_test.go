package common

import (
	"strings"
	"testing"
)

func TestGreeting(t *testing.T) {
	if got := Greeting("Ada"); !strings.Contains(got, "Ada") {
		t.Errorf("Greeting(Ada) = %q, want it to contain the name", got)
	}
}

func TestGreetingDefault(t *testing.T) {
	if got := Greeting(""); !strings.Contains(got, "world") {
		t.Errorf("Greeting(\"\") = %q, want default to mention world", got)
	}
}
