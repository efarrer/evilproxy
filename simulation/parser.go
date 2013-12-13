package simulation

import (
	"errors"
	"fmt"
)

func ConstructConnections(rule string) (Connection, Connection, error) {
	if rule == "" {
		c0, c1 := NewBasicConnections(
			NewBasicPipe(),
			NewBasicPipe())
		return c0, c1, nil
	}

	return nil, nil, errors.New(fmt.Sprintf("Unable to parse \"%v\"", rule))
}
