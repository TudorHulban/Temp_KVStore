package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRoundRobin(t *testing.T) {
	rr := NewRoundRobin(0, 3)

	require.EqualValues(t,
		1,
		rr.Next(),
	)
	require.EqualValues(t,
		2,
		rr.Next(),
	)
	require.EqualValues(t,
		3,
		rr.Next(),
	)
	require.EqualValues(t,
		0,
		rr.Next(),
	)
}
