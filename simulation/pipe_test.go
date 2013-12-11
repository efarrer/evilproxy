package simulation

import (
	"testing"
	"time"
)

/*
 * Note none of these tests will be automatically run
 * Implementors of the Pipe interface should call PerformPipeTests from their
 * unit tests to ensure their implementation is compliant.
 */

func testClosingAfterSendingStillDeliversPacket(pipeGenerator func() Pipe, t *testing.T) {
	pkt := Packet{}
	pipe := pipeGenerator()
	pipe.Send(&pkt)
	pipe.Close()
	rcvd, err := pipe.Recv()
	if err != nil {
		t.Fatalf("Got an unexpected error while receiving from pipe\n")
	}
	if &pkt != rcvd {
		t.Fatalf("Didn't get expected packet from pipe. Got %v expected %v", rcvd, pkt)
	}
}

func testPipeDeliversPacketsInOrder(pipeGenerator func() Pipe, t *testing.T) {
	pkt0 := &Packet{}
	pkt1 := &Packet{}
	pipe := pipeGenerator()
	defer pipe.Close()
	pipe.Send(pkt0)
	pipe.Send(pkt1)
	rcvd0, err := pipe.Recv()
	if err != nil {
		t.Fatalf("Got an unexpected error while receiving from pipe\n")
	}
	if pkt0 != rcvd0 {
		t.Fatalf("Didn't get first packet from pipe. Got %v expected %v",
			&rcvd0, &pkt0)
	}
	rcvd1, err := pipe.Recv()
	if err != nil {
		t.Fatalf("Got an unexpected error while receiving from pipe\n")
	}
	if pkt1 != rcvd1 {
		t.Fatalf("Didn't get second packet from pipe. Got %v expected %v",
			&rcvd1, &pkt1)
	}
}

func testSendingAfterCloseResultsInError(pipeGenerator func() Pipe, t *testing.T) {
	defer func() {
		recover()
	}()
	pkt := Packet{}
	pipe := pipeGenerator()
	pipe.Close()
	pipe.Send(&pkt)
	t.Fatalf("Expecting a panic for sending over closed pipe\n")
}

func testRecvHangsIfNoPacket(pipeGenerator func() Pipe, t *testing.T) {
	const delay = 100
	pkt := Packet{}
	pipe := pipeGenerator()
	go func() {
		<-time.After(time.Millisecond * delay)
		pipe.Send(&pkt)
	}()
	timer := StartTimer()
	rcvd, err := pipe.Recv()
	if err != nil {
		t.Fatalf("Got an unexpected error while receiving from pipe\n")
	}
	if FuzzyEquals(delay, timer.ElapsedMilliseconds(), 10) {
		t.Fatalf("Recv didn't block %v milliseconds expected %v milliseconds",
			timer.ElapsedMilliseconds(), delay)
	}
	if rcvd != &pkt {
		t.Fatalf("Didn't get expected packet from pipe. Got %v expected %v", rcvd, pkt)
	}
}

func testRecvFromClosedPipeResultsInNilPacket(pipeGenerator func() Pipe, t *testing.T) {
	pipe := pipeGenerator()
	pipe.Close()
	rcvd, err := pipe.Recv()
	if err == nil {
		t.Fatalf("Expecting an error receiving from closed pipe but got none.\n")
	}
	if rcvd != nil {
		t.Fatalf("Got a packet, but expected nil Got %v expected nil", rcvd)
	}
}

func testClosingAClosedPipePanics(pipeGenerator func() Pipe, t *testing.T) {
	defer func() {
		recover()
	}()
	pipe := pipeGenerator()
	pipe.Close()
	pipe.Close()
	t.Fatalf("Expected panic on double close")
}

func PerformPipeTests(pipeGenerator func() Pipe, t *testing.T) {

	testClosingAfterSendingStillDeliversPacket(pipeGenerator, t)
	testPipeDeliversPacketsInOrder(pipeGenerator, t)
	testSendingAfterCloseResultsInError(pipeGenerator, t)
	testRecvHangsIfNoPacket(pipeGenerator, t)
	testRecvFromClosedPipeResultsInNilPacket(pipeGenerator, t)
	testClosingAClosedPipePanics(pipeGenerator, t)
}
