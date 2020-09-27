package lmc

import "fmt"

// Instruction represents an instruction.
type Instruction struct {
	Label    string
	Mnemonic string
	Operand  string
	Opcode   int
}

// DefaultMnemonicMap maps mnemonics to their opcodes and default operands.
var DefaultMnemonicMap = map[string]Instruction{
	"ADD": {"", "ADD", "0", 1},
	"SUB": {"", "SUB", "0", 2},
	"STA": {"", "STA", "0", 3},
	"STO": {"", "STO", "0", 3},
	"LDA": {"", "LDA", "0", 5},
	"BRA": {"", "BRA", "0", 6},
	"BRZ": {"", "BRZ", "0", 7},
	"BRP": {"", "BRP", "0", 8},
	"INP": {"", "INP", "1", 9},
	"OUT": {"", "OUT", "2", 9},
	"DAT": {"", "DAT", "0", -1}, // DAT is special, doesn't have opcode.
	"HLT": {"", "HLT", "0", 0},
}

// Parser converts a stream of tokens into a list of Instructions.
type Parser struct {
	lexer *Lexer

	curToken  Token
	peekToken Token

	curInstruction Instruction
}

// NewParser returns a new parser from a lexer.
func NewParser(lexer *Lexer) *Parser {
	parser := &Parser{lexer: lexer}
	parser.readToken()
	parser.readToken()

	return parser
}

// readToken reads one more token from the input.
func (p *Parser) readToken() {
	p.curToken = p.peekToken
	p.peekToken = p.lexer.Next()
}

// Parse converts a stream of tokens into a list of Instructions.
func (p *Parser) Parse() ([]Instruction, error) {
	instructions := []Instruction{}

	for p.curToken.Type != EOF {
		switch p.curToken.Type {
		case ILLEGAL:
			return nil, fmt.Errorf("illegal token %s in input", p.curToken)

		case IDENT:
			if p.peekToken.Type == IDENT && DefaultMnemonicMap[p.peekToken.Literal].Mnemonic != "" {

				// This is a label before a mnemonic
				p.curInstruction.Label = p.curToken.Literal
				p.readToken()

			} else {

				// This is a mnemonic
				instruction := DefaultMnemonicMap[p.curToken.Literal]
				if instruction.Mnemonic == "" {
					// This mnemonic isn't valid/doesn't exist.
					return nil, fmt.Errorf("invalid mnemonic: %s", p.curToken.Literal)
				}

				// Set mnemonic
				p.curInstruction.Mnemonic = instruction.Mnemonic

				// Set opcdoe
				p.curInstruction.Opcode = instruction.Opcode

				// Set special case operands for INP, OUT.
				// Then set other operands.
				if instruction.Mnemonic == "INP" {
					p.curInstruction.Operand = "1"
				} else if instruction.Mnemonic == "OUT" {
					p.curInstruction.Operand = "2"
				} else if p.peekToken.Type == INT || p.peekToken.Type == IDENT {
					p.readToken()
					p.curInstruction.Operand = p.curToken.Literal
				}

				p.readToken()

			}

		case NEWLINE: // Use newlines as a cue to mean instruction has finished
			if p.curInstruction.Mnemonic == "" {
				p.readToken()
				continue
			}

			p.readToken()
			instructions = append(instructions, p.curInstruction)
			p.curInstruction = Instruction{}

		default:
			return nil, fmt.Errorf("unexpected token in input: %s", p.curToken)
		}
	}

	if p.curInstruction.Mnemonic != "" {
		instructions = append(instructions, p.curInstruction)
	}

	return instructions, nil
}
