package simulation

import (
	"testing"
)

func TestPipeBehaviorForBasicPipe(t *testing.T) {
	PerformPipeTests(func() Pipe { return NewBasicPipe() }, t)
}
