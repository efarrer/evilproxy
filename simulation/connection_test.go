package simulation

import (
	"testing"
	"time"
)

/*
 * Note none of these tests will be automatically run
 * Implementors of the connection interface should call PerformConnectionTests
 * from their unit tests to ensure their implementation is compliant.
 */

func testClosingAfterWritingStillDeliversPacket(
	connectionGenerator func() (Connection, Connection), t *testing.T) {
	pkt := Packet{}
	c0, c1 := connectionGenerator()
	c0.Write(&pkt)
	c0.Close()
	read, err := c1.Read()
	if err != nil {
		t.Fatalf("Got an unexpected error while reading from connection\n")
	}
	if &pkt != read {
		t.Fatalf("Didn't get expected packet from connection. Got %v expected %v", read, pkt)
	}
}

func testConnectionDeliversPacketsInOrder(
	connectionGenerator func() (Connection, Connection), t *testing.T) {
	pkt0 := &Packet{}
	pkt1 := &Packet{}
	c0, c1 := connectionGenerator()
	defer c0.Close()
	defer c1.Close()
	c0.Write(pkt0)
	c0.Write(pkt1)
	rcvd0, err := c1.Read()
	if err != nil {
		t.Fatalf("Got an unexpected error while receiving from connection\n")
	}
	if pkt0 != rcvd0 {
		t.Fatalf("Didn't get first packet from connection. Got %v expected %v",
			&rcvd0, &pkt0)
	}
	rcvd1, err := c1.Read()
	if err != nil {
		t.Fatalf("Got an unexpected error while receiving from connection\n")
	}
	if pkt1 != rcvd1 {
		t.Fatalf("Didn't get second packet from connection. Got %v expected %v",
			&rcvd1, &pkt1)
	}
}

func testWriteingAfterCloseResultsInError(
	connectionGenerator func() (Connection, Connection), t *testing.T) {
	defer func() {
		recover()
	}()
	pkt := Packet{}
	c0, _ := connectionGenerator()
	c0.Close()
	c0.Write(&pkt)
	t.Fatalf("Expecting a panic for writing over closed connection\n")
}

func testReadHangsIfNoPacket(
	connectionGenerator func() (Connection, Connection), t *testing.T) {
	const delay = 100
	pkt := Packet{}
	c0, c1 := connectionGenerator()
	go func() {
		<-time.After(time.Millisecond * delay)
		c0.Write(&pkt)
	}()
	timer := StartTimer()
	read, err := c1.Read()
	if err != nil {
		t.Fatalf("Got an unexpected error while reading from connection\n")
	}
	if FuzzyEquals(delay, timer.ElapsedMilliseconds(), 10) {
		t.Fatalf("Read didn't block %v milliseconds expected %v milliseconds",
			timer.ElapsedMilliseconds(), delay)
	}
	if read != &pkt {
		t.Fatalf("Didn't get expected packet from connection. Got %v expected %v", read, pkt)
	}
}

func testReadFromClosedPeerConnectionResultsInNilPacket(
	connectionGenerator func() (Connection, Connection), t *testing.T) {
	c0, c1 := connectionGenerator()
	c0.Close()
	read, err := c1.Read()
	if err == nil {
		t.Fatalf("Expecting an error reading from closed connection but got none.\n")
	}
	if read != nil {
		t.Fatalf("Got a packet, but expected nil Got %v expected nil", read)
	}
}

func testClosingAClosedConnectionPanics(
	connectionGenerator func() (Connection, Connection), t *testing.T) {
	defer func() {
		recover()
	}()
	c0, _ := connectionGenerator()
	c0.Close()
	c0.Close()
	t.Fatalf("Expected panic on double close")
}

func PerformConnectionTests(
	connectionGenerator func() (Connection, Connection), t *testing.T) {

	testClosingAfterWritingStillDeliversPacket(connectionGenerator, t)
	testConnectionDeliversPacketsInOrder(connectionGenerator, t)
	testWriteingAfterCloseResultsInError(connectionGenerator, t)
	testReadHangsIfNoPacket(connectionGenerator, t)
	testReadFromClosedPeerConnectionResultsInNilPacket(connectionGenerator, t)
	testClosingAClosedConnectionPanics(connectionGenerator, t)
}
