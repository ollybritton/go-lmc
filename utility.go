package lmc

import (
	"fmt"
	"strings"
)

// isWhitespace returns true if the character given is whitespace. It does not return true for newlines
// as these have to be handled by the lexer.
func isWhitespace(ch byte) bool {
	switch ch {
	case ' ', '\r', '\t':
		return true
	default:
		return false
	}
}

// isLetter returns true if the input character is a letter (a-z or A-Z).
func isLetter(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

// isIdentifier returns true if a string is a valid identifier.
func isIdentifier(s string) bool {
	bs := []byte(s)

	if !isLetter(bs[0]) {
		return false
	}

	for _, ch := range bs {
		if !isLetter(ch) && !isDigit(ch) {
			return false
		}
	}

	return true
}

// isDigit returns true if the input character is a number (0-9).
func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

// leftPadInt takes an integer and prepends zeros until it's the desired length.
func leftPadInt(n int, size int) string {
	s := fmt.Sprint(n)
	return strings.Repeat("0", size-len(s)) + s
}
