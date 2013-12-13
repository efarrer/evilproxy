package simulation

import (
	"testing"
)

func TestParsingBogusRuleReturnsError(t *testing.T) {
	cconn, sconn, err := ConstructConnections("bogus")

	if cconn != nil || sconn != nil || err == nil {
		t.Fatalf("Expecting error but got %v, %v\n", cconn, sconn)
	}
}

func TestParsingEmptyStringReturnsBasicConnectionsWithBasicPipes(t *testing.T) {
	cconn, sconn, err := ConstructConnections("")

	if err != nil {
		t.Fatalf("Expecting basicConnections of basicPipes got an error %v\n", err)
	}

	defer func() {
		if p := recover(); p != nil {
			t.Fatalf("Expecting basicConnections of basicPipes got something else.\n")
		}
	}()

	cc := cconn.(*basicConnection)
	sc := sconn.(*basicConnection)
	_ = cc.send.(*basicPipe)
	_ = cc.recv.(*basicPipe)
	_ = cc.close.(*basicPipe)
	_ = sc.send.(*basicPipe)
	_ = sc.recv.(*basicPipe)
	_ = sc.close.(*basicPipe)
}
