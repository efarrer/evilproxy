package evil_proxy

/*
 * A pipe is a thread-safe, bidirectional communication channel for
 * transmitting data. Packet's written with 'Write' will be available via a
 * peer's 'Read' method and visa-versa.
 */
type Connection interface {
	/*
	 * Queue the 'Packet' for writing.
	 * Panics if the connection is closed.
	 */
	Write(*Packet)

	/*
	 * Closes the connection.
	 */
	Close()

	/*
	 * Reads a 'Packet'.
	 * Blocks if a 'Packet' is not immediately available.
     * Returns error if the connection's peer is closed and all queued 'Packet's
     * have been read.
	 */
	Read() (*Packet, error)
}
