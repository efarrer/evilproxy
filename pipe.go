package evil_proxy

type Sender interface {
	// Queue the packet for sending
	// Panics if the sender is closed
	Send(*Packet)

	// Closes the Sender
	Close()
}

type Receiver interface {
	// Receives a packet
	// Waits for a packet if one is not immediately available
	Recv() (*Packet, error)
}

/*
 * A pipe is a unidirectional communication channel
 */
type Pipe interface {
	Sender
	Receiver
}
