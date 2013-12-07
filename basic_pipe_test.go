package evil_proxy

import (
	"testing"
)

func TestPipeBehaviorForBasicPipe(t *testing.T) {
	PerformPipeTests(func() Pipe { return NewBasicPipe() }, t)
}
