package parser

import (
	"errors"
	"fmt"

	"github.com/efarrer/evilproxy/connection"
	"github.com/efarrer/evilproxy/pipe"
)

func ConstructConnections(rule string) (connection.Connection, connection.Connection, error) {
	if rule == "" {
		c0, c1 := connection.NewBasicConnections(
			pipe.NewBasicPipe(),
			pipe.NewBasicPipe())
		return c0, c1, nil
	}

	return nil, nil, errors.New(fmt.Sprintf("Unable to parse \"%v\"", rule))
}
