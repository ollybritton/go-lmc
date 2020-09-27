package lmc_test

import (
	"testing"

	"github.com/ollybritton/go-lmc"
	"github.com/stretchr/testify/assert"
)

func TestComputerFromMailboxes(t *testing.T) {
	tests := []struct {
		input         []string
		inbox, outbox []int
		pc            int
		accumulator   int
	}{
		{
			[]string{"101", "010"},
			[]int{}, []int{},
			2, 10,
		},
	}

	for _, tc := range tests {
		mbs := lmc.NewInitialisedMailboxes(1, 2, tc.input)
		computer := lmc.NewComputerFromMailboxes(mbs, 1, 2)

		for _, val := range tc.inbox {
			computer.Inbox <- val
		}

		err := computer.Run()
		assert.NoError(t, err, "not expecting error executing program")

		assert.Equal(t, tc.pc, computer.ProgramCounter, "expect program counter to be correct")
		assert.Equal(t, tc.accumulator, computer.Accumulator, "expect accumulator to be correct")
	}
}
