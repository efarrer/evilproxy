package connection

import (
	"testing"

	"github.com/efarrer/evilproxy/pipe"
)

func TestConnectionBehaviorForBasicConnection(t *testing.T) {
	PerformConnectionTests(func() (Connection, Connection) {
		return NewBasicConnections(pipe.NewBasicPipe(), pipe.NewBasicPipe())
	}, t)
}
