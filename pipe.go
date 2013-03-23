package evil_proxy

type Sender interface {
	// Send the packet
	// Panics if the sender is closed
	Send(*Packet)

	// Closes the Sender
	Close()
}

type Receiver interface {
	Recv() (*Packet, error)
}

/*
 * A pipe is a unidirectional communication channel
 */
type Pipe interface {
	Sender
	Receiver
}
