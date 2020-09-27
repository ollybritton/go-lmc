package lmc_test

import (
	"testing"

	"github.com/ollybritton/go-lmc"
	"github.com/stretchr/testify/assert"
)

func TestParserValid(t *testing.T) {
	input := `INP
STA 10
ADD 5
label DAT 3
hello OUT
num1 DAT 0
INP
STA num1
num1 DAT 0`

	lexer := lmc.NewLexer(input)
	parser := lmc.NewParser(lexer)
	tests := []lmc.Instruction{
		{Label: "", Mnemonic: "INP", Operand: "1", Opcode: 9},
		{Label: "", Mnemonic: "STA", Operand: "10", Opcode: 3},
		{Label: "", Mnemonic: "ADD", Operand: "5", Opcode: 1},
		{Label: "label", Mnemonic: "DAT", Operand: "3", Opcode: -1},
		{Label: "hello", Mnemonic: "OUT", Operand: "2", Opcode: 9},
		{Label: "num1", Mnemonic: "DAT", Operand: "0", Opcode: -1},
		{Label: "", Mnemonic: "INP", Operand: "1", Opcode: 9},
		{Label: "", Mnemonic: "STA", Operand: "num1", Opcode: 3},
		{Label: "num1", Mnemonic: "DAT", Operand: "0", Opcode: -1},
	}

	instructions, err := parser.Parse()
	assert.Nil(t, err, "not expecting error when executing parser")

	for i, tc := range tests {
		instruction := instructions[i]

		assert.Equal(t, tc.Label, instruction.Label, "expect instruction label to be correct; want %s but got %s", tc.Label, instruction.Label)
		assert.Equal(t, tc.Mnemonic, instruction.Mnemonic, "expect instruction mnemonic to be correct; want %s but got %s", tc.Mnemonic, instruction.Mnemonic)
		assert.Equal(t, tc.Operand, instruction.Operand, "expect instruction operand to be correct; want %s but got %s", tc.Operand, instruction.Operand)
		assert.Equal(t, tc.Opcode, instruction.Opcode, "expect instruction opcode to be correct; want %d but got %d", tc.Opcode, instruction.Opcode)
	}

	final := lexer.Next()
	assert.Equal(t, lmc.EOF, final.Type, "expected last token to be EOF")

}

func TestParserInvalid(t *testing.T) {
	tests := []struct {
		name  string
		input string
		err   string
	}{
		{
			"double-label",
			"label label",
			"invalid mnemonic: label",
		},
		{
			"just-number",
			"10",
			"unexpected token in input: 10<INT>(line=0,col=0-0)",
		},
		{
			"invalid-char",
			"STA @",
			"illegal token @<ILLEGAL>(line=0,col=4-4) in input",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			lexer := lmc.NewLexer(tc.input)
			parser := lmc.NewParser(lexer)

			_, err := parser.Parse()
			assert.NotNil(t, err, "expecting error")
			assert.Equal(t, tc.err, err.Error(), "expected different error")
		})
	}

}
