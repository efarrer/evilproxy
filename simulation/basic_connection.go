package simulation

type basicConnection struct {
	send  Sender
	recv  Receiver
	close Closer
}

func (c *basicConnection) Write(p *Packet) {
	c.send.Send(p)
}

func (c *basicConnection) Close() {
	c.close.Close()
}

func (c *basicConnection) Read() (*Packet, error) {
	return c.recv.Recv()
}

/*
 * Constructs a pair of connections that use the given pipes to communicate.
 * The first connection will write using the first pipe's 'Send' method, and will
 * read using the second Pipe's 'Recv' method.
 */
func NewBasicConnections(p0, p1 Pipe) (Connection, Connection) {
	return &basicConnection{p0, p1, p0}, &basicConnection{p1, p0, p1}
}
