package parser

import (
	"testing"
)

func TestParsingBogusRuleReturnsError(t *testing.T) {
	cconn, sconn, err := ConstructConnections("bogus")

	if cconn != nil || sconn != nil || err == nil {
		t.Fatalf("Expecting error but got %v, %v\n", cconn, sconn)
	}
}

func TestParsingEmptyStringSucceedes(t *testing.T) {
	_, _, err := ConstructConnections("")

	if err != nil {
		t.Fatalf("Expecting basicConnections of basicPipes got an error %v\n", err)
	}
}
