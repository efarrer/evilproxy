package evil_proxy

/*
 * A socket is an endpoint of a bidirectional communication channel
 */
type Socket interface {
	Send([]byte)
	Recv() []byte
	Close()
}
