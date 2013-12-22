package simulation

import (
	"io"
)

/*
 * A connection is a thread-safe, bidirectional communication channel for
 * transmitting data. Packet's written with 'Write' will be available via a
 * peer's 'Read' method and visa-versa.
 */
type Connection interface {
	/*
	 * Queue the 'Packet' for writing.
	 */
	Write(*Packet) error

	/*
	 * Closes the connection.
	 */
	io.Closer

	/*
	 * Reads a 'Packet'.
	 * Blocks if a 'Packet' is not immediately available.
	 * Returns error if the connection's peer is closed and all queued 'Packet's
	 * have been read.
	 */
	Read() (*Packet, error)
}
