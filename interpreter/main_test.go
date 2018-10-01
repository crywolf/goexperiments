package main

import (
	"bytes"
	"strings"
	"testing"
)

func Test_interpret(t *testing.T) {

	input1 := `
# Ignore this line because it is a comment
COPY _Var1    10
COPY   _Var2 20

ADD _Var1 _Var2 # one could add a comment here
PRINT _Var1
PRINT _Var2
`

	output1 := "30\n20"

	tests := []struct {
		name    string
		input   string
		output  string
		wantErr bool
	}{
		{"correct commands", input1, output1, false},
		{"unknown command", "Cmd", "", true},
		{"ADD: opt1 does not start with underscore", "ADD a", "", true},
		{"ADD: opt2 is not integer or variable", "ADD _a rr", "", true},
		{"ADD: correct copying", "COPY _a -321\nPRINT _a", "-321", false},
		{"COPY: correct addition", "ADD _b 5\nPRINT _b", "5", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.input)
			writer := &bytes.Buffer{}
			errWriter := &bytes.Buffer{}

			interpret(reader, writer, errWriter)

			if err := errWriter.String(); (err != "") != tt.wantErr {
				t.Errorf("interpret() error = %v, want %v", err, tt.wantErr)
			}

			got := strings.TrimSpace(writer.String())
			want := tt.output
			if got != want {
				t.Errorf("interpret() = %v, want %v", got, want)
			}
		})
	}
}
