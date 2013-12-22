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
	err := pipe.Send(&pkt)
	UnexpectedError(err, "sending", t)
	err = pipe.Close()
	UnexpectedError(err, "closing", t)
	rcvd, err := pipe.Recv()
	UnexpectedError(err, "recving", t)
	if &pkt != rcvd {
		t.Fatalf("Didn't get expected packet from pipe. Got %v expected %v", rcvd, pkt)
	}
}

func testPipeDeliversPacketsInOrder(pipeGenerator func() Pipe, t *testing.T) {
	pkt0 := &Packet{}
	pkt1 := &Packet{}
	pipe := pipeGenerator()
	defer pipe.Close()
	err := pipe.Send(pkt0)
	UnexpectedError(err, "sending", t)
	err = pipe.Send(pkt1)
	UnexpectedError(err, "sending", t)
	rcvd0, err := pipe.Recv()
	UnexpectedError(err, "recving", t)
	if pkt0 != rcvd0 {
		t.Fatalf("Didn't get first packet from pipe. Got %v expected %v",
			&rcvd0, &pkt0)
	}
	rcvd1, err := pipe.Recv()
	UnexpectedError(err, "recving", t)
	if pkt1 != rcvd1 {
		t.Fatalf("Didn't get second packet from pipe. Got %v expected %v",
			&rcvd1, &pkt1)
	}
}

func testSendingAfterCloseResultsInError(pipeGenerator func() Pipe, t *testing.T) {
	pkt := Packet{}
	pipe := pipeGenerator()
	err := pipe.Close()
	UnexpectedError(err, "closing", t)
	err = pipe.Send(&pkt)
	if err == nil {
		t.Fatalf("Expecting error for sending over closed pipe\n")
	}
}

func testRecvHangsIfNoPacket(pipeGenerator func() Pipe, t *testing.T) {
	const delay = 100
	pkt := Packet{}
	pipe := pipeGenerator()
	go func() {
		<-time.After(time.Millisecond * delay)
		err := pipe.Send(&pkt)
		UnexpectedError(err, "recving", t)
	}()
	timer := StartTimer()
	rcvd, err := pipe.Recv()
	UnexpectedError(err, "recving", t)
	if FuzzyEquals(delay, timer.ElapsedMilliseconds(), 10) {
		t.Fatalf("Recv didn't block %v milliseconds expected %v milliseconds",
			timer.ElapsedMilliseconds(), delay)
	}
	if rcvd != &pkt {
		t.Fatalf("Didn't get expected packet from pipe. Got %v expected %v", rcvd, pkt)
	}
}

func testRecvFromClosedPipeResultsInNilPacketAndError(pipeGenerator func() Pipe, t *testing.T) {
	pipe := pipeGenerator()
	err := pipe.Close()
	UnexpectedError(err, "closing", t)
	rcvd, err := pipe.Recv()
	if rcvd != nil {
		t.Fatalf("Got a packet, but expected nil Got %v expected nil", rcvd)
	}
	if err == nil {
		t.Fatalf("Didn't get expected error when Recv'ing from closed pipe")
	}
}

func testClosingAClosedPipeFails(pipeGenerator func() Pipe, t *testing.T) {
	pipe := pipeGenerator()
	err := pipe.Close()
	UnexpectedError(err, "closing", t)
	err = pipe.Close()
	if err == nil {
		t.Fatalf("Expected error on double close")
	}
}

func PerformPipeTests(pipeGenerator func() Pipe, t *testing.T) {

	testClosingAfterSendingStillDeliversPacket(pipeGenerator, t)
	testPipeDeliversPacketsInOrder(pipeGenerator, t)
	testSendingAfterCloseResultsInError(pipeGenerator, t)
	testRecvHangsIfNoPacket(pipeGenerator, t)
	testRecvFromClosedPipeResultsInNilPacketAndError(pipeGenerator, t)
	testClosingAClosedPipeFails(pipeGenerator, t)
}
