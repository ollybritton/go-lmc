package lmc

import "fmt"

// Lexer is a lexer for LMC assembly.
// It translates some series of characters into tokens.
type Lexer struct {
	input string // Program input.
	ch    byte   // Current character under examination.

	position      int // Current position.
	positionNext  int // Position of the next character.
	positionStart int // Position of the current token under examination.

	line int // Current line.
	col  int // Current column in line.
}

// NewLexer returns a new, initialised instance of Lexer.
func NewLexer(input string) *Lexer {
	lexer := &Lexer{
		input: input,
		col:   -1, // Set to 0 when char is first read, zero-indexed.
	}

	lexer.readChar()
	return lexer
}

// readChar reads the next character in the input. If there's no more characters left
// to be read, it is set to the NUL character.
func (l *Lexer) readChar() {
	if l.positionNext >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.positionNext]
		l.col++
	}

	l.position = l.positionNext
	l.positionNext++
}

// peekChar returns the next character in the input. Like readChar, if there's no more characters
// left, then it returns the NUL character.
func (l *Lexer) peekChar() byte {
	if l.positionNext >= len(l.input) {
		return 0
	}

	return l.input[l.positionNext]
}

// skipWhitespace consumes and discards all whitespace until the next non-whitespace character.
func (l *Lexer) skipWhitespace() {
	for isWhitespace(l.ch) {
		l.readChar()
	}
}

// skipComment will skip any comments over any number of lines.
func (l *Lexer) skipComment() {
	if l.ch != '/' && l.peekChar() != '/' {
		fmt.Println("returning")
		return
	}

	for l.ch != '\n' {
		l.readChar()
	}
}

// readIdentifier reads any sequence of valid identifier characters and returns a string and the start and end index
// of the identifier relative to the current line.
func (l *Lexer) readIdentifier() (string, int, int) {
	l.positionStart = l.position
	colStart := l.col

	for isLetter(l.ch) || (isDigit(l.ch) && l.position != l.positionStart) {
		l.readChar()
	}

	return l.input[l.positionStart:l.position], colStart, l.col - 1
}

// readInteger reads a number and returns it's value as a string, along with the start and end index of the integer
// relative to the current line.
func (l *Lexer) readInteger() (string, int, int) {
	l.positionStart = l.position
	colStart := l.col

	for isDigit(l.ch) {
		l.readChar()
	}

	if colStart == l.col {
		return l.input[l.positionStart:l.position], colStart, l.col
	}

	return l.input[l.positionStart:l.position], colStart, l.col - 1
}

// Next returns the next token in the input.
func (l *Lexer) Next() Token {
	if isWhitespace(l.ch) {
		l.skipWhitespace()
	}

	if l.ch == '/' {
		l.skipComment()
	}

	switch {
	case isDigit(l.ch):
		lit, start, end := l.readInteger()
		return NewToken(INT, lit, l.line, start, end)

	case isLetter(l.ch):
		lit, start, end := l.readIdentifier()
		return NewToken(IDENT, lit, l.line, start, end)

	case l.ch == 0:
		return NewToken(EOF, "", l.line, l.col, l.col)

	case l.ch == '\n':
		tok := NewToken(NEWLINE, "\n", l.line, l.col, l.col)

		l.readChar()

		l.line++
		l.col = 0

		return tok
	}

	curr := l.ch
	l.readChar()

	return NewToken(ILLEGAL, string(curr), l.line, l.col, l.col)
}
