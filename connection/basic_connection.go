package connection

import (
	"io"

	"github.com/efarrer/evilproxy/packet"
	"github.com/efarrer/evilproxy/pipe"
)

type basicConnection struct {
	send  pipe.Sender
	recv  pipe.Receiver
	close io.Closer
}

func (c *basicConnection) Write(p *packet.Packet) error {
	return c.send.Send(p)
}

func (c *basicConnection) Close() error {
	return c.close.Close()
}

func (c *basicConnection) Read() (*packet.Packet, error) {
	return c.recv.Recv()
}

/*
 * Constructs a pair of connections that use the given pipes to communicate.
 * The first connection will write using the first pipe's 'Send' method, and will
 * read using the second Pipe's 'Recv' method.
 */
func NewBasicConnections(p0, p1 pipe.Pipe) (Connection, Connection) {
	return &basicConnection{p0, p1, p0}, &basicConnection{p1, p0, p1}
}
