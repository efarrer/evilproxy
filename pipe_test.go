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

type Seconds int64

func (t timer) elapsedSeconds() Seconds {
	milliseconds := time.Now().Sub(time.Time(t)) / time.Millisecond

	// Round to the nearest second
	seconds := ((milliseconds + 500) / 1000)

	return Seconds(seconds)
}

func TestLatentPipeDelaysPacketsIfGivenDelay(t *testing.T) {
	pkt := Packet{}
	pipe := NewLatentPipe(time.Millisecond*1000)
	timer := startTimer()
	pipe.Send(&pkt)
	rcvd := pipe.Recv()
	if 1 != timer.elapsedSeconds() {
		t.Fatalf("Latent pipe didn't express expected latency. Took %v seconds expected %v seconds",
			timer.elapsedSeconds(), 1)
	}
	if &pkt != rcvd {
		t.Fatalf("Didn't get expected packet from latent pipe. Got %v expected %v", rcvd, pkt)
	}
}

func TestLatentPipeWontDelayIfNoDelayGiven(t *testing.T) {
	pkt := Packet{}
	pipe := NewLatentPipe(time.Millisecond*0)
	timer := startTimer()
	pipe.Send(&pkt)
	rcvd := pipe.Recv()
	if 0 != timer.elapsedSeconds() {
		t.Fatalf("Latent pipe expressed latency. Took %v seconds expected %v seconds", timer.elapsedSeconds(), 0)
	}
	if &pkt != rcvd {
		t.Fatalf("Didn't get expected packet from latent pipe. Got %v expected %v", rcvd, pkt)
	}
}

func TestClosingAfterSendingStillResultsInDeliveredPacket(t *testing.T) {
	pkt := Packet{}
	pipe := NewLatentPipe(time.Millisecond*0)
	pipe.Send(&pkt)
	pipe.Close()
	rcvd := pipe.Recv()
	if &pkt != rcvd {
		t.Fatalf("Didn't get expected packet from latent pipe. Got %v expected %v", rcvd, pkt)
	}
}

func TestClosingWithoutResultsInNilPacket(t *testing.T) {
	pipe := NewLatentPipe(time.Millisecond*0)
	pipe.Close()
	rcvd := pipe.Recv()
	if rcvd != nil {
		t.Fatalf("Got a packet, but expected nil Got %v expected nil", rcvd)
	}
}
