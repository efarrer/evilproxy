package simulation

import (
	"io"
)

type Sender interface {
	Send(*Packet) error
}

type Receiver interface {
	Recv() (*Packet, error)
}

/*
 * A 'Pipe' is a thread-safe, unidirectional communication channel for
 * transmitting data. Packet's sent with 'Send' will be available on the pipe's
 * own 'Recv' interface.
 */
type Pipe interface {
	/*
	 * Queue the 'Packet' for sending.
	 */
	Sender

	/*
	 * Closes the pipe.
	 */
	io.Closer

	/*
	 * Receives a 'Packet'.
	 * Blocks if a 'Packet' is not immediately available.
	 * Returns error if the pipe is closed and all queued 'Packet's have been
	 * received.
	 */
	Receiver
}
