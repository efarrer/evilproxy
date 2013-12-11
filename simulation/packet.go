package simulation

type Flags byte

const (
	None Flags = 0
	Syn  Flags = 1
	Ack  Flags = 2
	Fin  Flags = 4
)

type Packet struct {
	Flags      Flags
	Seq        int64
	Ack        int64
	WindowSize int64
	Payload    []byte
}
