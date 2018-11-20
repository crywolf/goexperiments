package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
)

func main() {
	commands := [][]string{
		{"echo", "Hello, \nWorld"},
		{"make", "all"},
		{"ls", "-l"},
		{"git", "status"},
	}

	done := make(chan error, len(commands))

	for _, cmd := range commands {
		c := exec.Command(cmd[0], cmd[1:]...)

		prefixedOutputWriter := newPrefixingWriter(cmd[0], os.Stdout)
		c.Stdout = prefixedOutputWriter
		c.Stderr = prefixedOutputWriter

		go func() {
			done <- c.Run()
		}()
	}

	for i := 0; i < len(commands); i++ {
		<-done
	}
}

func newPrefixingWriter(prefix string, output io.Writer) io.Writer {
	reader, writer := io.Pipe()

	scanner := bufio.NewScanner(reader)

	go func() {
		defer writer.Close()

		for scanner.Scan() {
			// Write the prefix
			fmt.Fprintf(output, "[%s] ", prefix)
			// Copy the line
			output.Write(scanner.Bytes())
			// Re-add a new line (scanner removes it)
			fmt.Fprint(output, "\n")
		}
	}()

	return writer
}
