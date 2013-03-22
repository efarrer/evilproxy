package evil_proxy

type Sender interface {
	Send(*Packet)
}

type Receiver interface {
	Recv() *Packet
}

/*
 * A pipe is a unidirectional communication channel
 */
type Pipe interface {
	Sender
	Receiver
	Close()
}
