package cmd

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestCommand(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		input   string
		wantW   string
		wantErr bool
	}{
		{
			name:    "simple input",
			args:    []string{""},
			input:   "{ \"message\": \"Praise the Sun\" }",
			wantW:   "Praise the Sun\n",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := []string{""}
			RootCmd.SetArgs(args)

			gotW := executeCommandWithInput(t, tt.input)

			if gotW != tt.wantW {
				t.Errorf("executeCommandWithInput gotW = '%v', want '%v'", gotW, tt.wantW)
			}
		})
	}
}

func executeCommandWithInput(t *testing.T, input string) string {
	inputBytes := []byte(input)
	r := prepareStdIn(t, inputBytes)
	readStdOut, w := prepareStdOut(t)
	stdin := os.Stdin
	stdout := os.Stdout
	defer func() {
		os.Stdin = stdin
		os.Stdout = stdout
	}()
	os.Stdin = r
	os.Stdout = w

	err := RootCmd.Execute()
	if err != nil {
		t.Fatal(err)
	}
	err = w.Close()
	if err != nil {
		t.Fatal(err)
	}

	output, err := ioutil.ReadAll(readStdOut)
	if err != nil {
		t.Fatal(err)
	}
	return string(output)
}

func prepareStdIn(t *testing.T, input []byte) *os.File {
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	_, err = w.Write(input)
	if err != nil {
		t.Error(err)
	}
	err = w.Close()
	if err != nil {
		t.Error(err)
	}
	return r
}

func prepareStdOut(t *testing.T) (*os.File, *os.File) {
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	return r, w
}