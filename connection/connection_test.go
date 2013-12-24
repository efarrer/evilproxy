package connection

import (
	"errors"
	"evilproxy/packet"
	"evilproxy/testing_utils"
	"testing"
	"time"
)

/*
 * Note none of these tests will be automatically run
 * Implementors of the connection interface should call PerformConnectionTests
 * from their unit tests to ensure their implementation is compliant.
 */

func testClosingAfterWritingStillDeliversPacket(
	connectionGenerator func() (Connection, Connection), t *testing.T) {
	pkt := packet.Packet{}
	c0, c1 := connectionGenerator()
	err := c0.Write(&pkt)
	testing_utils.UnexpectedError(err, "writing", t)
	err = c0.Close()
	testing_utils.UnexpectedError(err, "closing", t)
	read, err := c1.Read()
	testing_utils.UnexpectedError(err, "reading", t)
	if &pkt != read {
		t.Fatalf("Didn't get expected packet from connection. Got %v expected %v", read, pkt)
	}
}

func testConnectionDeliversPacketsInOrder(
	connectionGenerator func() (Connection, Connection), t *testing.T) {
	pkt0 := &packet.Packet{}
	pkt1 := &packet.Packet{}
	c0, c1 := connectionGenerator()
	defer c0.Close()
	defer c1.Close()
	err := c0.Write(pkt0)
	testing_utils.UnexpectedError(err, "writing", t)
	err = c0.Write(pkt1)
	testing_utils.UnexpectedError(err, "writing", t)
	rcvd0, err := c1.Read()
	testing_utils.UnexpectedError(err, "reading", t)
	if pkt0 != rcvd0 {
		t.Fatalf("Didn't get first packet from connection. Got %v expected %v",
			&rcvd0, &pkt0)
	}
	rcvd1, err := c1.Read()
	testing_utils.UnexpectedError(err, "reading", t)
	if pkt1 != rcvd1 {
		t.Fatalf("Didn't get second packet from connection. Got %v expected %v",
			&rcvd1, &pkt1)
	}
}

func testWriteingAfterCloseResultsInError(
	connectionGenerator func() (Connection, Connection), t *testing.T) {
	pkt := packet.Packet{}
	c0, _ := connectionGenerator()
	err := c0.Close()
	testing_utils.UnexpectedError(err, "closing", t)
	err = c0.Write(&pkt)
	if err == nil {
		t.Fatalf("Expecting error for writing over a closed connection\n")
	}
}

func testReadHangsIfNoPacket(
	connectionGenerator func() (Connection, Connection), t *testing.T) {
	const delay = 100
	pkt := packet.Packet{}
	c0, c1 := connectionGenerator()
	go func() {
		<-time.After(time.Millisecond * delay)
		c0.Write(&pkt)
	}()
	timer := testing_utils.StartTimer()
	read, err := c1.Read()
	testing_utils.UnexpectedError(err, "reading", t)
	if testing_utils.FuzzyEquals(delay, timer.ElapsedMilliseconds(), 10) {
		t.Fatalf("Read didn't block %v milliseconds expected %v milliseconds",
			timer.ElapsedMilliseconds(), delay)
	}
	if read != &pkt {
		t.Fatalf("Didn't get expected packet from connection. Got %v expected %v", read, pkt)
	}
}

func testReadFromClosedPeerConnectionResultsInNilPacketAndError(
	connectionGenerator func() (Connection, Connection), t *testing.T) {
	c0, c1 := connectionGenerator()
	c0.Close()
	read, err := c1.Read()
	if read != nil {
		t.Fatalf("Got a packet, but expected nil Got %v expected nil", read)
	}
	if err == nil {
		t.Fatalf("Didn't get expected error when Read'ing from closed connection")
	}
}

func testClosingAClosedConnectionFails(
	connectionGenerator func() (Connection, Connection), t *testing.T) {
	c0, _ := connectionGenerator()
	err := c0.Close()
	testing_utils.UnexpectedError(err, "closing", t)
	err = c0.Close()
	if err == nil {
		t.Fatalf("Expected error on double close")
	}
}

func PerformConnectionTests(
	connectionGenerator func() (Connection, Connection), t *testing.T) {

	testClosingAfterWritingStillDeliversPacket(connectionGenerator, t)
	testConnectionDeliversPacketsInOrder(connectionGenerator, t)
	testWriteingAfterCloseResultsInError(connectionGenerator, t)
	testReadHangsIfNoPacket(connectionGenerator, t)
	testReadFromClosedPeerConnectionResultsInNilPacketAndError(connectionGenerator, t)
	testClosingAClosedConnectionFails(connectionGenerator, t)
}

type packetReader struct {
	err    error
	packet *packet.Packet
}

func (p *packetReader) Read() (*packet.Packet, error) {
	if p.packet != nil {
		pkt := p.packet
		p.packet = nil
		return pkt, nil
	}
	return nil, p.err
}

func TestConnectionReaderAdaptorReturnsErrorOnError(t *testing.T) {
	pktReader := &packetReader{errors.New("Some error"), nil}

	reader := ConnectionReaderAdaptor(pktReader)

	buffer := make([]byte, 10)
	n, err := reader.Read(buffer)
	if n != 0 {
		t.Fatalf("Expected read of 0 got %v\n", n)
	}
	if err == nil {
		t.Fatalf("Expected read error, but no error was returned\n")
	}
}

func TestConnectionReaderAdaptorReturnsDataFromPacket(t *testing.T) {
	const DataSize = 10
	pkt := &packet.Packet{}
	pkt.Payload = make([]byte, DataSize)
	pktReader := &packetReader{errors.New("Some error"), pkt}

	buffer := make([]byte, DataSize)
	reader := ConnectionReaderAdaptor(pktReader)

	n, err := reader.Read(buffer)
	if n != DataSize {
		t.Fatalf("Expected read of %v got %v\n", DataSize, n)
	}
	if err != nil {
		t.Fatalf("Got an unexpected error %v\n", err)
	}
	n, err = reader.Read(buffer)
	if n != 0 {
		t.Fatalf("Expected read of 0 got %v\n", n)
	}
	if err == nil {
		t.Fatalf("Expected read error, but no error was returned\n")
	}
}

func TestConnectionReaderAdaptorCanReadMultipleTimesFromLargePacket(t *testing.T) {
	const DataSize = 10
	const PacketSizeMultiplier = 5
	pkt := &packet.Packet{}
	pkt.Payload = make([]byte, DataSize*PacketSizeMultiplier)
	pktReader := &packetReader{errors.New("Some error"), pkt}

	buffer := make([]byte, DataSize)
	reader := ConnectionReaderAdaptor(pktReader)

	for i := 0; i != PacketSizeMultiplier; i++ {
		n, err := reader.Read(buffer)
		if n != DataSize {
			t.Fatalf("Expected read of %v got %v\n", DataSize, n)
		}
		if err != nil {
			t.Fatalf("Got an unexpected error %v\n", err)
		}
	}
	n, err := reader.Read(buffer)
	if n != 0 {
		t.Fatalf("Expected read of 0 got %v\n", n)
	}
	if err == nil {
		t.Fatalf("Expected read error, but no error was returned\n")
	}
}
