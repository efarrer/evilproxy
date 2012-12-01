package evil_proxy

/*
 * A pipe is a unidirectional communication channel
 */
type Pipe interface {
	Send(*Packet)
	Recv() *Packet
	Close()
}
