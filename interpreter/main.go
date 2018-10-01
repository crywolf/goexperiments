package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func main() {
	interpret(os.Stdin, os.Stdout, os.Stderr)
}

func interpret(reader io.Reader, writer io.Writer, errWriter io.Writer) {
	variables := make(map[string]int)

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		tokens := strings.Fields(line)

		var cmdName string
		var op1 string
		var op2 int
		var err error
		var errOccurred = false

		for i, token := range tokens {
			// skip comments
			if strings.HasPrefix(token, "#") {
				break
			}
			// first is command name
			if i == 0 {
				switch token {
				case "COPY", "ADD", "PRINT":
					cmdName = token
				default:
					errOccurred = true
					break
				}
			}
			// operand 1
			if i == 1 {
				if !strings.HasPrefix(token, "_") {
					errOccurred = true
					break
				}
				op1 = token
			}
			// operand 2
			if i == 2 {
				if strings.HasPrefix(token, "_") {
					res, varExists := variables[token]
					if !varExists {
						errOccurred = true
						break
					}
					op2 = res
				} else {
					op2, err = strconv.Atoi(token)
					if err != nil {
						errOccurred = true
						break
					}
				}
			}
		}
		if errOccurred {
			printError(errWriter)
			continue
		}

		switch cmdName {
		case "COPY":
			variables[op1] = op2
		case "ADD":
			variables[op1] += op2
		case "PRINT":
			res, varExists := variables[op1]
			if !varExists {
				printError(errWriter)
				continue
			}
			fmt.Fprintf(writer, "%d\n", res)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(errWriter, "error reading input: %v", err)
	}
}

func printError(writer io.Writer) {
	fmt.Fprintln(writer, "Error!")
}
