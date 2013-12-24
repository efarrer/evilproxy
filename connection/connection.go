package connection

import (
	"evilproxy/packet"
	"io"
)

type PacketReader interface {
	Read() (*packet.Packet, error)
}

/*
 * A connection is a thread-safe, bidirectional communication channel for
 * transmitting data. Packet's written with 'Write' will be available via a
 * peer's 'Read' method and visa-versa.
 */
type Connection interface {
	/*
	 * Queue the 'Packet' for writing.
	 */
	Write(*packet.Packet) error

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
	PacketReader
}

type readerAdaptor struct {
	currentPacket *packet.Packet
	pktReader     PacketReader
}

func (ra *readerAdaptor) Read(p []byte) (int, error) {
	if ra.currentPacket == nil {
		pkt, err := ra.pktReader.Read()
		if err != nil {
			return 0, err
		}
		ra.currentPacket = pkt
	}

	n := copy(p, ra.currentPacket.Payload)
	if n == len(ra.currentPacket.Payload) {
		ra.currentPacket = nil
	} else {
		ra.currentPacket.Payload = ra.currentPacket.Payload[n:]
	}
	return n, nil
}

func ConnectionReaderAdaptor(pktReader PacketReader) io.Reader {
	return &readerAdaptor{nil, pktReader}
}
