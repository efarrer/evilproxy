package evil_proxy

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

type clientSocket struct {
	pipe Pipe
}

func (s *clientSocket) Send(data []byte) {
	// TODO implement this
}

func (s *clientSocket) Recv() []byte {
	// TODO implement this
	return nil
}

func (s *clientSocket) Close() {
	// TODO implement this
}

/*
 * Constructs a new client for communicating over a pipe
 */
func NewClient(pipe Pipe) Socket {
	// Construct initial syn packet
	syn := Packet{Syn, rand.Int63(), 0, 0, []byte{}}
	pipe.Send(&syn)

	return &clientSocket{pipe}
}
