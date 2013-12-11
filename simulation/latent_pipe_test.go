package simulation

import (
	"testing"
	"time"
)

func TestPipeBehaviorForLatentPipe(t *testing.T) {
	PerformPipeTests(func() Pipe { return NewLatentPipe(NewBasicPipe(), time.Millisecond*0) }, t)
}

func TestLatentPipeDelaysPackets(t *testing.T) {
	const delay = 100
	pkt := Packet{}
	pipe := NewLatentPipe(NewBasicPipe(), time.Millisecond*delay)
	timer := StartTimer()
	pipe.Send(&pkt)
	rcvd, err := pipe.Recv()
	if err != nil {
		t.Fatalf("Got an unexpected error while receiving from pipe\n")
	}
	if FuzzyEquals(delay, timer.ElapsedMilliseconds(), 10) {
		t.Fatalf("Latent pipe didn't express expected latency. Took %v milliseconds expected %v milliseconds",
			timer.ElapsedMilliseconds(), delay)
	}
	if &pkt != rcvd {
		t.Fatalf("Didn't get expected packet from latent pipe. Got %v expected %v", rcvd, pkt)
	}
}

func TestLatentPipeWontDelayIfNoDelay(t *testing.T) {
	pkt := Packet{}
	pipe := NewLatentPipe(NewBasicPipe(), time.Millisecond*0)
	timer := StartTimer()
	pipe.Send(&pkt)
	rcvd, err := pipe.Recv()
	if err != nil {
		t.Fatalf("Got an unexpected error while receiving from pipe\n")
	}
	if FuzzyEquals(0, timer.ElapsedMilliseconds(), 10) {
		t.Fatalf("Latent pipe expressed latency. Took %v milliseconds expected %v milliseconds", timer.ElapsedMilliseconds(), 0)
	}
	if &pkt != rcvd {
		t.Fatalf("Didn't get expected packet from latent pipe. Got %v expected %v", rcvd, pkt)
	}
}
