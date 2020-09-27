package lmc_test

import (
	"testing"

	"github.com/ollybritton/go-lmc"
	"github.com/stretchr/testify/assert"
)

func TestAssemble(t *testing.T) {
	tests := []struct {
		input          string
		insize, opsize int
		output         []string
	}{
		{
			`ADD 10`,
			1, 2,
			[]string{"110", "000"},
		},
		{
			`INP
			 STA 10
			 DAT 05`,
			1, 2,
			[]string{"901", "310", "005"},
		},
		{
			`INP
			 STA num1
			 num1 DAT 0`,
			1, 2,
			[]string{"901", "302", "000"},
		},
		{
			`		INP
					STA VALUE
					LDA ZERO
					STA TRINUM
					STA N
			LOOP    LDA TRINUM
					SUB VALUE
					BRP ENDLOOP
					LDA N
					ADD ONE
					STA N
					ADD TRINUM
					STA TRINUM
					BRA LOOP
			ENDLOOP LDA VALUE
					SUB TRINUM
					BRZ EQUAL
					LDA ZERO
					OUT
					BRA DONE
			EQUAL   LDA N
					OUT
			DONE    HLT
			VALUE   DAT
			TRINUM  DAT
			N       DAT
			ZERO    DAT 000
			ONE     DAT 001`,
			1, 2,
			[]string{
				"901", "323", "526", "324", "325", "524", "223", "814", "525", "127",
				"325", "124", "324", "605", "523", "224", "720", "526", "902", "622",
				"525", "902", "000", "000", "000", "000", "000", "001", "000", "000",
			},
		},
		{
			`ADD 10`,
			2, 3,
			[]string{"01010"},
		},
		{
			`INP
			 STA 10
			 DAT 05`,
			2, 3,
			[]string{"09001", "03010", "00005"},
		},
		{
			`		INP
					STA VALUE
					LDA ZERO
					STA TRINUM
					STA N
			LOOP    LDA TRINUM
					SUB VALUE
					BRP ENDLOOP
					LDA N
					ADD ONE
					STA N
					ADD TRINUM
					STA TRINUM
					BRA LOOP
			ENDLOOP LDA VALUE
					SUB TRINUM
					BRZ EQUAL
					LDA ZERO
					OUT
					BRA DONE
			EQUAL   LDA N
					OUT
			DONE    HLT
			VALUE   DAT
			TRINUM  DAT
			N       DAT
			ZERO    DAT 000
			ONE     DAT 001`,
			2, 3,
			[]string{
				"09001", "03023", "05026", "03024", "03025", "05024", "02023", "08014", "05025", "01027",
				"03025", "01024", "03024", "06005", "05023", "02024", "07020", "05026", "09002", "06022",
				"05025", "09002", "00000", "00000", "00000", "00000", "00000", "00001", "00000", "00000",
			},
		},
	}

	for _, tc := range tests {
		lexer := lmc.NewLexer(tc.input)
		parser := lmc.NewParser(lexer)
		instructions, err := parser.Parse()
		assert.Nil(t, err, "not expecting error when executing parser")

		mailboxes := lmc.Assemble(instructions, tc.insize, tc.opsize)

		for i, val := range tc.output {
			got, err := mailboxes.Get(i)
			assert.NoError(t, err, "not expecting error accessing mailbox")
			assert.Equal(t, val, got, "expecting mailbox %d to be %s, got %s", i, val, got)
		}
	}
}
