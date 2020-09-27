package lmc

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// Mailboxes represents a memory for the Little Man Computer.
type Mailboxes struct {
	mem []string
}

// Get attempts to retrieve what is in mailbox N, indexed from 0.
func (m *Mailboxes) Get(n int) (string, error) {
	if n > len(m.mem)-1 {
		return "", ErrInvalidMemory{n}
	}

	if n < 0 {
		return "", ErrInvalidMemory{n}
	}

	return m.mem[n], nil
}

// Set attempts to set what is in mailbox N, indexed from 0.
func (m *Mailboxes) Set(n int, val string) error {
	if n > len(m.mem)-1 {
		return ErrInvalidMemory{n}
	}

	if n < 0 {
		return ErrInvalidMemory{n}
	}

	m.mem[n] = val
	return nil
}

// NewMailboxes returns a new Mailboxes instance with the size specified.
func NewMailboxes(inSize, opSize int) *Mailboxes {
	size := int(math.Pow(10, float64(inSize+opSize)))

	mb := &Mailboxes{
		mem: make([]string, size),
	}

	def := strings.Repeat("0", inSize+opSize)
	for i := 0; i < size; i++ {
		mb.Set(i, def)
	}

	return mb
}

// NewInitialisedMailboxes returns a new Mailboes instance with the size specified and the values given set.
func NewInitialisedMailboxes(inSize, opSize int, init []string) *Mailboxes {
	mb := NewMailboxes(inSize, opSize)

	for i, val := range init {
		mb.Set(i, val)
	}

	return mb
}

// Computer represents a Little Man Computer that supports a variable n
type Computer struct {
	Mailboxes      *Mailboxes
	ProgramCounter int
	Accumulator    int

	InstructionSize int
	OperandSize     int

	Messages chan Msg
	Step     chan struct{}
	Inbox    chan int
}

// Status represents a status of the computer.
type Status string

// Status definitions
const (
	NeedInput Status = "NeedInput"
	NeedStep  Status = "NeedStup"
	Log       Status = "Log"
	Done      Status = "Done"
	Output    Status = "Output"
)

// Msg is a message sent to the user of a Little Man Computer.
type Msg struct {
	Status Status // What kind of message this is.
	Val    string // Optional value, used for passingl ogs.
}

// NewComputerFromMailboxes returns a new computer that will execute the instructions in the mailboxes.
// The in/out boxes by default have a 100 item buffer.
func NewComputerFromMailboxes(mailboxes *Mailboxes, inSize, opSize int) *Computer {
	return &Computer{
		Mailboxes:       mailboxes,
		InstructionSize: inSize,
		OperandSize:     opSize,
		Messages:        make(chan Msg),
		Step:            make(chan struct{}),
		Inbox:           make(chan int),
	}
}

// NewComputerFromCode returns a new computer whose mailboxes are loaded with the source code given.
func NewComputerFromCode(code string, inSize, opSize int) (*Computer, error) {
	lexer := NewLexer(code)
	parser := NewParser(lexer)

	instructions, err := parser.Parse()
	if err != nil {
		return nil, err
	}

	mailboxes := Assemble(instructions, inSize, opSize)
	return NewComputerFromMailboxes(mailboxes, inSize, opSize), nil
}

// Run starts the computer, it will continue until it hits a HLT instruction.
func (c *Computer) Run() error {
	c.Messages <- Msg{Log, "Little Man warming up..."}

	for i := 0; true; i++ {
		if i != 0 {
			c.Messages <- Msg{NeedStep, ""}
			<-c.Step
		}

		c.Messages <- Msg{Log, fmt.Sprintf("Getting instruction/operand at address %d", c.ProgramCounter)}
		memNum, err := c.Mailboxes.Get(c.ProgramCounter)
		if err != nil {
			return err
		}

		memStr := fmt.Sprint(memNum)
		memStr = strings.Repeat("0", len(memStr)-(c.InstructionSize+c.OperandSize)) + memStr

		instructionStr := memStr[0:c.InstructionSize]
		operandStr := memStr[c.InstructionSize:]

		c.Messages <- Msg{Log, fmt.Sprintf("Instruction code: %s, Operand: %s", instructionStr, operandStr)}

		instruction, err := strconv.Atoi(instructionStr)
		if err != nil {
			c.Messages <- Msg{Done, ""}
			return err
		}

		operand, err := strconv.Atoi(operandStr)
		if err != nil {
			c.Messages <- Msg{Done, ""}
			return err
		}

		switch instruction {
		case 1: // ADD
			c.Messages <- Msg{Log, fmt.Sprintf("ADD; adding what is at address %d to accumulator", operand)}
			val, err := c.Mailboxes.Get(operand)
			if err != nil {
				c.Messages <- Msg{Done, ""}
				return err
			}

			num, err := strconv.Atoi(val)
			if err != nil {
				c.Messages <- Msg{Done, ""}
				return err
			}

			c.Accumulator += num
			c.Messages <- Msg{Log, fmt.Sprintf("ADD; added %d to accumulator, new value %d", num, c.Accumulator)}

		case 2: // SUB
			c.Messages <- Msg{Log, fmt.Sprintf("SUB; subtracting what is at address %d from accumulator", operand)}
			val, err := c.Mailboxes.Get(operand)
			if err != nil {
				c.Messages <- Msg{Done, ""}
				return err
			}

			num, err := strconv.Atoi(val)
			if err != nil {
				c.Messages <- Msg{Done, ""}
				return err
			}

			c.Accumulator -= num
			c.Messages <- Msg{Log, fmt.Sprintf("SUB; subtracted %d from accumulator, new value %d", num, c.Accumulator)}

		case 3: // STA
			c.Messages <- Msg{Log, fmt.Sprintf("STA; storing accumulator %d at address %d", c.Accumulator, operand)}
			err := c.Mailboxes.Set(operand, leftPadInt(c.Accumulator, c.InstructionSize+c.OperandSize))
			if err != nil {
				c.Messages <- Msg{Done, ""}
				return err
			}

		case 5: // LDA
			c.Messages <- Msg{Log, fmt.Sprintf("LDA; loading what is at address %d into accumulator", operand)}
			val, err := c.Mailboxes.Get(operand)
			if err != nil {
				c.Messages <- Msg{Done, ""}
				return err
			}

			num, err := strconv.Atoi(val)
			if err != nil {
				c.Messages <- Msg{Done, ""}
				return err
			}

			c.Accumulator = num
			c.Messages <- Msg{Log, fmt.Sprintf("LDA; set accumulator to %d", c.Accumulator)}

		case 6: // BRA
			c.Messages <- Msg{Log, fmt.Sprintf("BRA; setting program counter to %d", operand)}
			c.ProgramCounter = operand
			continue

		case 7: // BRZ
			c.Messages <- Msg{Log, fmt.Sprintf("BRZ; setting program counter to %d if accumulator is zero", operand)}
			if c.Accumulator == 0 {
				c.ProgramCounter = operand
				c.Messages <- Msg{Log, fmt.Sprintf("BRZ; accumlator IS zero, program counter now %d", operand)}
				continue
			}

		case 8: // BRP
			c.Messages <- Msg{Log, fmt.Sprintf("BRP; setting program counter to %d if accumulator is positive", operand)}
			if c.Accumulator >= 0 {
				c.ProgramCounter = operand
				c.Messages <- Msg{Log, fmt.Sprintf("BRZ; accumlator IS positive (%d), program counter now %d", c.Accumulator, operand)}
				continue
			}

		case 9: // INP/OUT
			if operand == 1 {
				c.Messages <- Msg{Log, "INP; Need input from user"}
				c.Messages <- Msg{NeedInput, ""}

				val := <-c.Inbox
				c.Messages <- Msg{Log, fmt.Sprintf("INP; Recieved input %d from user, set accumlator", val)}
				c.Accumulator = val
			} else if operand == 2 {
				c.Messages <- Msg{Log, "OUT; Outputting accumulator contents"}
				c.Messages <- Msg{Output, fmt.Sprint(c.Accumulator)}
			}
		}

		if operand == 0 && instruction == 0 {
			c.Messages <- Msg{Log, "HLT; We're done here!"}
			c.Messages <- Msg{Done, ""}
			return nil
		}

		c.ProgramCounter++
	}

	return nil
}
