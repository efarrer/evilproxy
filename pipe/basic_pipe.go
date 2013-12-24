package pipe

import (
	"container/list"
	"errors"
	"evilproxy/packet"
)

/*
 * A basicPipe is the simplest pipe that satisfies the expected pipe behaviro
 */
type basicPipe struct {
	inputChan  chan *packet.Packet
	outputChan chan *packet.Packet
	closed     bool
}

/*
 * Send a packet over a basic pipe
 */
func (bp *basicPipe) Send(p *packet.Packet) error {
	if bp.closed {
		return errors.New("Sending on a closed basic pipe.\n")
	}
	bp.inputChan <- p
	return nil
}

/*
 * Receive a packet from the basic pipe
 */
func (bp *basicPipe) Recv() (*packet.Packet, error) {
	pkt, ok := <-bp.outputChan
	if !ok {
		return nil, errors.New("Receiver is closed.")
	}
	return pkt, nil
}

/*
 * Close the basic pipe
 */
func (bp *basicPipe) Close() error {
	if bp.closed {
		return errors.New("Closing a closed basic pipe.\n")
	}
	close(bp.inputChan)
	bp.closed = true
	return nil
}

/*
 * Constructs a new basic pipe
 */
func NewBasicPipe() Pipe {
	bp := &basicPipe{make(chan *packet.Packet), make(chan *packet.Packet), false}

	go func() {
		var shutdown = false

		// The packets that have arrived and are ready to be recv'd
		arrived := list.New()
		var arrived_head *packet.Packet = nil

		// Holds either bp.outputChan or nil
		// if outputChan is nil then it will never be selected
		// Should be set if arrived_head is set and nil if arrived_head is nil
		var outputChan chan *packet.Packet = nil

		for {
			// If we've been shutdown and the arrived queue is empty then we can
			// close the outputChan and exit
			if shutdown &&
				arrived_head == nil && arrived.Len() == 0 {
				close(bp.outputChan)
				return
			}

			select {
			// We got a new packet to queue
			case input, ok := <-bp.inputChan:
				// We've been shutdown, but we need to wait for packets in the
				// arrived queued data to be delivered, so just set a flag for
				// now.
				if !ok {
					shutdown = true
					continue
				}

				// Push this packet onto the arrived list
				if arrived_head == nil {
					arrived_head = input
					outputChan = bp.outputChan
				} else {
					arrived.PushBack(input)
				}

			case outputChan <- arrived_head:
				if 0 == arrived.Len() {
					arrived_head = nil
					outputChan = nil
				} else {
					elm := arrived.Front()
					arrived.Remove(elm)
					arrived_head = elm.Value.(*packet.Packet)
					outputChan = bp.outputChan
				}

			}
		}
	}()

	return bp
}
