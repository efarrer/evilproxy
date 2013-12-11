package simulation

type Sender interface {
	Send(*Packet)
}

type Receiver interface {
	Recv() (*Packet, error)
}

type Closer interface {
	Close()
}

/*
 * A 'Pipe' is a thread-safe, unidirectional communication channel for
 * transmitting data. Packet's sent with 'Send' will be available on the pipe's
 * own 'Recv' interface.
 */
type Pipe interface {
	/*
	 * Queue the 'Packet' for sending.
	 * Panics if the pipe is closed.
	 */
	Sender

	/*
	 * Closes the pipe.
	 */
	Closer

	/*
	 * Receives a 'Packet'.
	 * Blocks if a 'Packet' is not immediately available.
	 * Returns error if the pipe is closed and all queued 'Packet's have been
	 * received.
	 */
	Receiver
}
