package evil_proxy

type Sender interface {
	/*
	 * Queue the packet for sending.
	 * Panics if the sender is closed.
	 */
	Send(*Packet)

	// Closes the Sender
	Close()
}

type Receiver interface {
	/* Receives a packet.
		 * Waits for a packet if one is not immediately available.
	     * Returns error if sender has closed.
	*/
	Recv() (*Packet, error)
}

/*
 * A pipe is a thread-safe, unidirectional communication channel for
 * transmitting packets data sent with the Sender interface will be available on
 * the pipe's own Receiver interface.
 */
type Pipe interface {
	Sender
	Receiver
}
