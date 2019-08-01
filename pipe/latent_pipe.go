package pipe

import (
	"container/list"
	"errors"
	"time"

	"github.com/efarrer/evilproxy/packet"
)

/*
 * A latentPipe is a pipe that simulates latency
 */
type latentPipe struct {
	inputChan chan *latentPacket
	basePipe  Pipe
	latency   time.Duration
	closed    bool
}

/*
 * A container for packets that contain the arrival time
 */
type latentPacket struct {
	packet      *packet.Packet
	arrivalTime time.Time
}

/*
 * Send a packet over a latent pipe
 */
func (lp *latentPipe) Send(p *packet.Packet) error {
	if lp.closed {
		return errors.New("Sending on a closed latent pipe.\n")
	}
	lp.inputChan <- &latentPacket{p, time.Now().Add(lp.latency)}
	return nil
}

/*
 * Receive a packet from the latent pipe
 */
func (lp *latentPipe) Recv() (*packet.Packet, error) {
	return lp.basePipe.Recv()
}

/*
 * Close the latent pipe
 */
func (lp *latentPipe) Close() error {
	if lp.closed {
		return errors.New("Closing a closed latent pipe.\n")
	}
	close(lp.inputChan)
	lp.closed = true
	return nil
}

/*
 * Constructs a new latent pipe, with the given latency
 */
func NewLatentPipe(p Pipe, latency time.Duration) Pipe {
	lp := &latentPipe{make(chan *latentPacket), p, latency, false}

	go func() {
		var shutdown = false
		// The next packet on the queue to be sent
		var timer <-chan time.Time = nil
		restartTimer := func(lp *latentPacket) {
			timer = time.After(lp.arrivalTime.Sub(time.Now()))
		}

		// The packets that are in transit over the latent pipe
		intransit := list.New()
		var intransit_head *latentPacket = nil

		// The packets that have arrived and are ready to be recv'd
		arrived := list.New()
		var arrived_head *packet.Packet = nil

		for {
			// If we've been shutdown and the arrived queue is empty then we can
			// close the basePipe and exit
			if shutdown &&
				arrived_head == nil && arrived.Len() == 0 &&
				intransit_head == nil && intransit.Len() == 0 {
				lp.basePipe.Close()
				return
			}

			select {
			// We got a new packet to queue
			case input, ok := <-lp.inputChan:
				// We've been shutdown so make sure the other channel is also
				// closed then exit
				if !ok {
					shutdown = true
					continue
				}

				// If this is the first packet then start the timer
				if nil == intransit_head {
					intransit_head = input
					restartTimer(input)
				} else { // Push this packet on the end of the intransit list
					intransit.PushBack(input)
				}

			// The head packet is ready to deliver
			case <-timer:
				// Push this packet to the base pipe
				lp.basePipe.Send(intransit_head.packet)

				// Since we no longer have any packets in transit disable the
				// timer
				if intransit.Len() == 0 {
					timer = nil
					intransit_head = nil
				} else {
					elm := intransit.Front()
					intransit.Remove(elm)
					intransit_head = elm.Value.(*latentPacket)
					restartTimer(intransit_head)
				}

			}
		}
	}()

	return lp
}
