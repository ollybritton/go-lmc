package lmc_test

import (
	"testing"

	"github.com/ollybritton/go-lmc"
	"github.com/stretchr/testify/assert"
)

func TestLexerNext(t *testing.T) {
	input := `INP
STA 10
ADD 5
// Comment on one line
BZP 20 // Comment at end of line
   SUB 3
// Multiple comments
// Multiple comments
label ADD 2`

	lexer := lmc.NewLexer(input)
	tests := []lmc.Token{
		{Type: lmc.IDENT, Literal: "INP", Line: 0, StartCol: 0, EndCol: 2},
		{Type: lmc.NEWLINE, Literal: "\n", Line: 0, StartCol: 3, EndCol: 3},
		{Type: lmc.IDENT, Literal: "STA", Line: 1, StartCol: 0, EndCol: 2},
		{Type: lmc.INT, Literal: "10", Line: 1, StartCol: 4, EndCol: 5},
		{Type: lmc.NEWLINE, Literal: "\n", Line: 1, StartCol: 6, EndCol: 6},
		{Type: lmc.IDENT, Literal: "ADD", Line: 2, StartCol: 0, EndCol: 2},
		{Type: lmc.INT, Literal: "5", Line: 2, StartCol: 4, EndCol: 4},
		{Type: lmc.NEWLINE, Literal: "\n", Line: 2, StartCol: 5, EndCol: 5},
		{Type: lmc.NEWLINE, Literal: "\n", Line: 3, StartCol: 22, EndCol: 22},
		{Type: lmc.IDENT, Literal: "BZP", Line: 4, StartCol: 0, EndCol: 2},
		{Type: lmc.INT, Literal: "20", Line: 4, StartCol: 4, EndCol: 5},
		{Type: lmc.NEWLINE, Literal: "\n", Line: 4, StartCol: 32, EndCol: 32},
		{Type: lmc.IDENT, Literal: "SUB", Line: 5, StartCol: 3, EndCol: 5},
		{Type: lmc.INT, Literal: "3", Line: 5, StartCol: 7, EndCol: 7},
		{Type: lmc.NEWLINE, Literal: "\n", Line: 5, StartCol: 8, EndCol: 8},
		{Type: lmc.NEWLINE, Literal: "\n", Line: 6, StartCol: 20, EndCol: 20},
		{Type: lmc.NEWLINE, Literal: "\n", Line: 7, StartCol: 20, EndCol: 20},
		{Type: lmc.IDENT, Literal: "label", Line: 8, StartCol: 0, EndCol: 4},
		{Type: lmc.IDENT, Literal: "ADD", Line: 8, StartCol: 6, EndCol: 8},
		{Type: lmc.INT, Literal: "2", Line: 8, StartCol: 10, EndCol: 10},
	}

	for _, tc := range tests {
		tok := lexer.Next()

		assert.Equal(t, tc.Type, tok.Type, "expect token type to be correct; want %s but got %s", tc.String(), tok.String())
		assert.Equal(t, tc.Literal, tok.Literal, "expect literal value to be correct; want %s but got %s", tc.String(), tok.String())
		assert.Equal(t, tc.Line, tok.Line, "expect token line num to be correct; want %s but got %s", tc.String(), tok.String())
		assert.Equal(t, tc.StartCol, tok.StartCol, "expect token start column to be correct; want %s but got %s", tc.String(), tok.String())
		assert.Equal(t, tc.EndCol, tok.EndCol, "expect token end column to be correct; want %s but got %s", tc.String(), tok.String())
	}

	final := lexer.Next()
	assert.Equal(t, lmc.EOF, final.Type, "expected last token to be EOF")

}
