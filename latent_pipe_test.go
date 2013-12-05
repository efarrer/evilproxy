package evil_proxy

import (
	"testing"
	"time"
)

/*
 * A simple timer for seeing how long an operation lasts
 */
type timer time.Time

func startTimer() timer {
	return timer(time.Now())
}

func (t timer) elapsedMilliseconds() time.Duration {
	return time.Now().Sub(time.Time(t)) * time.Millisecond
}

func fuzzyEquals(a, b, delta time.Duration) bool {
	diff := a - b
	if diff < 0 {
		diff *= -1
	}
	return diff < delta
}

func TestLatentPipeDelaysPackets(t *testing.T) {
	const delay = 100
	pkt := Packet{}
	pipe := NewLatentPipe(time.Millisecond * delay)
	timer := startTimer()
	pipe.Send(&pkt)
	rcvd, err := pipe.Recv()
	if err != nil {
		t.Fatalf("Got an unexpected error while receiving from pipe\n")
	}
	if fuzzyEquals(delay, timer.elapsedMilliseconds(), 10) {
		t.Fatalf("Latent pipe didn't express expected latency. Took %v milliseconds expected %v milliseconds",
			timer.elapsedMilliseconds(), delay)
	}
	if &pkt != rcvd {
		t.Fatalf("Didn't get expected packet from latent pipe. Got %v expected %v", rcvd, pkt)
	}
}

func TestLatentPipeWontDelayIfNoDelay(t *testing.T) {
	pkt := Packet{}
	pipe := NewLatentPipe(time.Millisecond * 0)
	timer := startTimer()
	pipe.Send(&pkt)
	rcvd, err := pipe.Recv()
	if err != nil {
		t.Fatalf("Got an unexpected error while receiving from pipe\n")
	}
	if fuzzyEquals(0, timer.elapsedMilliseconds(), 10) {
		t.Fatalf("Latent pipe expressed latency. Took %v milliseconds expected %v milliseconds", timer.elapsedMilliseconds(), 0)
	}
	if &pkt != rcvd {
		t.Fatalf("Didn't get expected packet from latent pipe. Got %v expected %v", rcvd, pkt)
	}
}

func TestClosingAfterSendingStillDeliversPacket(t *testing.T) {
	pkt := Packet{}
	pipe := NewLatentPipe(time.Millisecond * 0)
	pipe.Send(&pkt)
	pipe.Close()
	rcvd, err := pipe.Recv()
	if err != nil {
		t.Fatalf("Got an unexpected error while receiving from pipe\n")
	}
	if &pkt != rcvd {
		t.Fatalf("Didn't get expected packet from latent pipe. Got %v expected %v", rcvd, pkt)
	}
}

func TestSendingAfterCloseResultsInError(t *testing.T) {
	defer func() {
		recover()
	}()
	pkt := Packet{}
	pipe := NewLatentPipe(time.Millisecond * 0)
	pipe.Close()
	pipe.Send(&pkt)
	t.Fatalf("Expecting a panic for sending over closed pipe\n")
}

func TestRecvHangsIfNoPacket(t *testing.T) {
	const delay = 100
	pkt := Packet{}
	pipe := NewLatentPipe(time.Millisecond * 0)
	go func() {
		<-time.After(time.Millisecond * delay)
		pipe.Send(&pkt)
	}()
	timer := startTimer()
	rcvd, err := pipe.Recv()
	if err != nil {
		t.Fatalf("Got an unexpected error while receiving from pipe\n")
	}
	if fuzzyEquals(delay, timer.elapsedMilliseconds(), 10) {
		t.Fatalf("Recv didn't block %v milliseconds expected %v milliseconds",
			timer.elapsedMilliseconds(), delay)
	}
	if rcvd != &pkt {
		t.Fatalf("Didn't get expected packet from latent pipe. Got %v expected %v", rcvd, pkt)
	}
}

func TestRecvFromClosedPipeResultsInNilPacket(t *testing.T) {
	pipe := NewLatentPipe(time.Millisecond * 0)
	pipe.Close()
	rcvd, err := pipe.Recv()
	if err == nil {
		t.Fatalf("Expecting an error receiving from closed pipe but got none.\n")
	}
	if rcvd != nil {
		t.Fatalf("Got a packet, but expected nil Got %v expected nil", rcvd)
	}
}
