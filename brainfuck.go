package main

import (
	"bytes"
	"time"
	"fmt"
)

// Struct of the context for a brainfuck machine.
type BrainfuckContext struct {
	// The size of the ram in cells.
	RamSize uint
	// The timeout in milliseconds.
	Timeout int64

	// The input buffer.
	input *bytes.Buffer
	// The output buffer.
	output *bytes.Buffer

	// The code to execute.
	code []rune

	// The ram cells.
	ram []byte
	// The cell pointer.
	ptr uint
	// The program counter.
	pc uint
}

// Initialize the brainfuck context with code and an input.
func (bf *BrainfuckContext) Init(code string, input string) {
	if (bf.RamSize == 0) {
		bf.RamSize = 256
	}

	if (bf.Timeout == 0) {
		bf.Timeout = 10000
	}

	bf.input = bytes.NewBufferString(input)
	bf.output = bytes.NewBufferString("")

	bf.code = []rune(code)

	bf.ram = make([]byte, bf.RamSize)
	bf.ptr = 0
	bf.pc = 0
}

// Execute the code in this context.
func (bf *BrainfuckContext) Exec() error {
	bf.output.Reset()
	start := time.Now().UnixNano()

	for c := bf.fetch(); c != 0; c = bf.fetch() {
		switch c {
		case '+': bf.inc()
		case '-': bf.dec()
		case '<': bf.left()
		case '>': bf.right()
		case '.': bf.write()
		case ',': bf.read()
		case '[': bf.loopOpen()
		case ']': bf.loopClose()
		}

		now := time.Now().UnixNano()
		if (now - start > bf.Timeout * 1e6) {
			return fmt.Errorf("timeout")
		}

		bf.pc++
	}

	return nil
}

// Get the current 'instruction'.
func (bf *BrainfuckContext) fetch() rune {
	if bf.pc < uint(len(bf.code)) {
		return bf.code[bf.pc]
	}

	return 0
}

// Increase the value of the cell currently pointed to.
func (bf *BrainfuckContext) inc() {
	bf.ram[bf.ptr]++
}

// Decrease the value of the cell currently pointed to.
func (bf *BrainfuckContext) dec() {
	bf.ram[bf.ptr]--
}

// Move the pointer left.
func (bf *BrainfuckContext) left() {
	bf.ptr = (bf.ptr + bf.RamSize - 1) % bf.RamSize
}

// Move the pointer right.
func (bf *BrainfuckContext) right() {
	bf.ptr = (bf.ptr + 1) % bf.RamSize
}

// Write a value to the output.
func (bf *BrainfuckContext) write() {
	if bf.ram[bf.ptr] < 128 {
		bf.output.WriteByte(bf.ram[bf.ptr])
	}
}

// Read a value from the input to the cell currently pointed to.
func (bf *BrainfuckContext) read() {
	bf.ram[bf.ptr], _ = bf.input.ReadByte()
}

// If the value currently pointed to is 0, jump to the matching ]
func (bf *BrainfuckContext) loopOpen() {
	if bf.ram[bf.ptr] != 0 {
		return
	}

	depth := 0

	for i := bf.pc; i < uint(len(bf.code)); i++ {
		if bf.code[i] == '[' {
			depth++
		} else if bf.code[i] == ']' {
			depth--
		}

		if depth == 0 {
			bf.pc = i
			break
		}
	}
}

// If the value currently pointed to is not 0, jump to the matching [
func (bf *BrainfuckContext) loopClose() {
	if bf.ram[bf.ptr] == 0 {
		return
	}

	depth := 0

	for i := bf.pc; i >= 0; i-- {
		if bf.code[i] == ']' {
			depth++
		} else if bf.code[i] == '[' {
			depth--
		}

		if depth == 0 {
			bf.pc = i
			break
		}
	}
}

// Get the output of the program as a string
func (bf *BrainfuckContext) Output() string {
	return bf.output.String()
}