package lmc

import "fmt"

// TokenType represents the type of a token, such as an identifier, label or integer.
type TokenType string

// Definitions of token Types.
const (
	ILLEGAL TokenType = "ILLEGAL"
	EOF     TokenType = "EOF"

	NEWLINE TokenType = "NEWLINE"

	IDENT TokenType = "IDENT"
	INT   TokenType = "INT"
)

// Token represents a small, easily categorizable chunk of text within the assembly.
type Token struct {
	Type    TokenType
	Literal string

	Line     int
	StartCol int
	EndCol   int
}

// NewToken returns a new instance of token.
func NewToken(tokenType TokenType, lit string, line, startCol, endCol int) Token {
	return Token{
		Type:     tokenType,
		Literal:  lit,
		Line:     line,
		StartCol: startCol,
		EndCol:   endCol,
	}
}

// String returns a string representation of a token.
// The format is LITERAL<TYPE>(line=LINE_NUM,col=START-END)
func (t Token) String() string {
	return fmt.Sprintf("%s<%s>(line=%d,col=%d-%d)", t.Literal, t.Type, t.Line, t.StartCol, t.EndCol)
}
