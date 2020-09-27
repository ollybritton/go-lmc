package lmc

import (
	"fmt"
	"strings"
)

// Assemble takes in a list of instructions and outputs mailboxes with those instructions loaded.
func Assemble(instructions []Instruction, opcodeSize, operandSize int) *Mailboxes {
	mailboxes := NewMailboxes(opcodeSize, operandSize)

	labelMap := make(map[string]string)

	for i, instruction := range instructions {
		if instruction.Label != "" {
			labelMap[instruction.Label] = leftPadInt(i, (opcodeSize + operandSize))
		}
	}

	for i, instruction := range instructions {
		if instruction.Mnemonic == "DAT" {
			if instruction.Operand != "" && isIdentifier(instruction.Operand) {
				mailboxes.Set(i, labelMap[instruction.Operand]) // TODO: check if label is valid
			} else {
				// TODO: check if number isn't valid
				mailboxes.Set(i, strings.Repeat("0", (operandSize+opcodeSize)-len(instruction.Operand))+instruction.Operand)
			}
		} else {
			opcode := fmt.Sprint(instruction.Opcode)
			operand := fmt.Sprint(instruction.Operand)

			if instruction.Operand != "" && isIdentifier(operand) {
				operand = fmt.Sprint(labelMap[instruction.Operand])[opcodeSize:] // TODO: check if label is valid
			}

			operand = strings.Repeat("0", operandSize-len(operand)) + operand
			opcode = strings.Repeat("0", opcodeSize-len(opcode)) + opcode
			mailboxes.Set(i, opcode+operand)
		}
	}

	return mailboxes
}
