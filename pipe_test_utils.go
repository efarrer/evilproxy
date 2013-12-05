package evil_proxy

import (
	"testing"
	"time"
)

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
		t.Fatalf("Didn't get expected packet from latent pipe. Got %v expected %v", rcvd, pkt)
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

func PerformPipeTests(pipeGenerator func() Pipe, t *testing.T) {

	testClosingAfterSendingStillDeliversPacket(pipeGenerator, t)
	testSendingAfterCloseResultsInError(pipeGenerator, t)
	testRecvHangsIfNoPacket(pipeGenerator, t)
	testRecvFromClosedPipeResultsInNilPacket(pipeGenerator, t)
}
