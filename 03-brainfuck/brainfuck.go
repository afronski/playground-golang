package main

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

type Instruction struct {
	operator uint16
	operand  uint16
}

const (
	MOVE_RIGHT = iota
	MOVE_LEFT
	INCREMENT
	DECREMENT
	PUT_CHAR
	GET_CHAR
	LOOP
	WHILE
)

const TAPE_SIZE int = 65535

func compile(input string) (program []Instruction, err error) {
	var PROGRAM_COUNTER uint16 = 0
	var JUMP_COUNTER uint16 = 0

	JUMP_STACK := make([]uint16, 0)

	for _, c := range input {
		switch c {
		case '>':
			program = append(program, Instruction{MOVE_RIGHT, 0})

		case '<':
			program = append(program, Instruction{MOVE_LEFT, 0})

		case '+':
			program = append(program, Instruction{INCREMENT, 0})

		case '-':
			program = append(program, Instruction{DECREMENT, 0})

		case '.':
			program = append(program, Instruction{PUT_CHAR, 0})

		case ',':
			program = append(program, Instruction{GET_CHAR, 0})

		case '[':
			program = append(program, Instruction{LOOP, 0})
			JUMP_STACK = append(JUMP_STACK, PROGRAM_COUNTER)

		case ']':
			if len(JUMP_STACK) == 0 {
				return nil, errors.New("Compilation error.")
			}

			JUMP_COUNTER = JUMP_STACK[len(JUMP_STACK)-1]
			JUMP_STACK = JUMP_STACK[:len(JUMP_STACK)-1]

			program = append(program, Instruction{WHILE, JUMP_COUNTER})
			program[JUMP_COUNTER].operand = PROGRAM_COUNTER

		default:
			PROGRAM_COUNTER--
		}

		PROGRAM_COUNTER++
	}

	if len(JUMP_STACK) != 0 {
		return nil, errors.New("Compilation error.")
	}

	return
}

func execute(program []Instruction) {
	TAPE := make([]int16, TAPE_SIZE)
	var TAPE_POINTER uint16 = 0

	reader := bufio.NewReader(os.Stdin)

	for PROGRAM_COUNTER := 0; PROGRAM_COUNTER < len(program); PROGRAM_COUNTER++ {
		switch program[PROGRAM_COUNTER].operator {
		case MOVE_RIGHT:
			TAPE_POINTER++

		case MOVE_LEFT:
			TAPE_POINTER--

		case INCREMENT:
			TAPE[TAPE_POINTER]++

		case DECREMENT:
			TAPE[TAPE_POINTER]--

		case PUT_CHAR:
			fmt.Printf("%c", TAPE[TAPE_POINTER])

		case GET_CHAR:
			readCharacter, _ := reader.ReadByte()
			TAPE[TAPE_POINTER] = int16(readCharacter)

		case LOOP:
			if TAPE[TAPE_POINTER] == 0 {
				PROGRAM_COUNTER = int(program[PROGRAM_COUNTER].operand)
			}

		case WHILE:
			if TAPE[TAPE_POINTER] > 0 {
				PROGRAM_COUNTER = int(program[PROGRAM_COUNTER].operand)
			}

		default:
			panic("Unknown operator.")
		}
	}
}

func main() {
	args := os.Args

	if len(args) != 2 {
		fmt.Printf("Usage: %s filename\n", args[0])
		return
	}

	filename := args[1]
	fileContents, err := ioutil.ReadFile(filename)

	if err != nil {
		fmt.Printf("Error reading %s\n", filename)
		return
	}

	program, err := compile(string(fileContents))

	if err != nil {
		fmt.Println(err)
		return
	}

	execute(program)
}
