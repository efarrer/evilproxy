package evil_proxy

import (
	"container/list"
	"errors"
	"time"
)

/*
 * A latentPipe is a pipe that simulates latency
 */
type latentPipe struct {
	inputChan  chan *latentPacket
	outputChan chan *Packet
	latency    time.Duration
}

/*
 * A container for packets that contain the arrival time
 */
type latentPacket struct {
	packet      *Packet
	arrivalTime time.Time
}

/*
 * Send a packet over a latent pipe
 */
func (lp latentPipe) Send(p *Packet) {
	lp.inputChan <- &latentPacket{p, time.Now().Add(lp.latency)}
}

/*
 * Receive a packet from the latent pipe
 */
func (lp latentPipe) Recv() (*Packet, error) {
	pkt, ok := <-lp.outputChan
	if !ok {
		return nil, errors.New("Receiver is closed.")
	}
	return pkt, nil
}

/*
 * Close the latent pipe
 */
func (lp latentPipe) Close() {
	close(lp.inputChan)
}

/*
 * Constructs a new latent pipe, with the given latency
 */
func NewLatentPipe(latency time.Duration) Pipe {
	lp := latentPipe{make(chan *latentPacket), make(chan *Packet), latency}

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
		var arrived_head *Packet = nil

		// Holds either lp.outputChan or nil
		// if outputChan is nil then it will never be selected
		// Should be set if arrived_head is set and nil if arrived_head is nil
		var outputChan chan *Packet = nil

		for {
			// If we've been shutdown and the arrived queue is empty then we can
			// close the outputChan and exit
			if shutdown &&
				arrived_head == nil && arrived.Len() == 0 &&
				intransit_head == nil && intransit.Len() == 0 {
				close(lp.outputChan)
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
				if nil == intransit.Front() {
					intransit_head = input
					restartTimer(input)
				} else { // Push this packet on the end of the intransit list
					intransit.PushBack(input)
				}

			// The head packet is ready to deliver
			case <-timer:
				// Push this packet onto the arrived list
				if arrived.Len() == 0 {
					arrived_head = intransit_head.packet
					outputChan = lp.outputChan
				} else {
					arrived.PushBack(intransit_head.packet)
				}

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

			case outputChan <- arrived_head:
				if 0 == arrived.Len() {
					arrived_head = nil
					outputChan = nil
				} else {
					elm := arrived.Front()
					arrived.Remove(elm)
					arrived_head = elm.Value.(*Packet)
					outputChan = lp.outputChan
				}

			}
		}
	}()

	return lp
}
