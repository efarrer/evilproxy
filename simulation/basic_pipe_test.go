package simulation

import (
	"testing"
)

func TestPipeBehaviorForBasicPipe(t *testing.T) {
	PerformPipeTests(func() Pipe { return NewBasicPipe() }, t)
}

func BenchmarkSendingPacketsOverPipe(b *testing.B) {
	pkt := &Packet{}
	p := NewBasicPipe()
	defer p.Close()
	for i := 0; i < b.N; i++ {
		p.Send(pkt)
		p.Recv()
	}
}
