package evil_proxy

import (
	"testing"
)

func TestNewClientSendsSynPacket(t *testing.T) {
	pipe := NewLatentPipe(0)
	NewClient(pipe)
	pkt, err := pipe.Recv()
	if err != nil {
		t.Fatalf("Pipe shouldn't be closed")
	}
	if (pkt.Flags & Syn) != Syn {
		t.Fatalf("Client didn't send a syn packet")
	}
	if pkt.Ack != 0 {
		t.Fatalf("Client didn't send a zero ack field")
	}
	if pkt.Seq == 0 {
		t.Fatalf("Client sent a zero seq field")
	}
}
