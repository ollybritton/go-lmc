package lmc

import "fmt"

// ErrArgumentLen is an error that occurs when more or less than the desired amount of arguments are
// specified.
type ErrArgumentLen struct {
	Mnemonic string
	Operands []int
	Wanted   int
}

// Error returns the error string for ErrArgumentLen.
func (e ErrArgumentLen) Error() string {
	return fmt.Sprintf("%s wants %d arguments, got %d", e.Mnemonic, e.Wanted, len(e.Operands))
}

// ErrInvalidMemory occurs when a mailbox is asked for that does not exist.
type ErrInvalidMemory struct {
	Attempted int
}

// Error returns the error string for ErrInvalidMemory.
func (e ErrInvalidMemory) Error() string {
	return fmt.Sprintf("invalid memory access; attempted to get mailbox %d", e.Attempted)
}
