package cmd

import (
	"fmt"
	"github.com/mdreem/json2line/configuration"
	"github.com/mdreem/json2line/configuration/load"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestCommand(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		input   string
		wantW   string
		wantErr string
	}{
		{
			name:    "simple input",
			args:    []string{""},
			input:   "{ \"message\": \"Praise the Sun\" }",
			wantW:   "Praise the Sun\n",
			wantErr: "",
		},
		{
			name:    "check empty configuration",
			args:    []string{"configure", "-s"},
			input:   "",
			wantW:   "Templates:\n\nReplacements:\n\nConfiguration:\n",
			wantErr: "",
		},
		{
			name:    "check buffer_size command",
			args:    []string{"-b", "5"},
			input:   "{ \"message\": \"Praise the Sun\" }",
			wantW:   "exited with code 1\n",
			wantErr: "could not parse line: bufio.Scanner: token too long\n",
		},
		{
			name:    "test ad-hoc templates",
			args:    []string{"-o", "MY_TEMPLATE {{ .message }}.", "-b", "4096"},
			input:   "{ \"message\": \"Praise the Sun\" }",
			wantW:   "MY_TEMPLATE Praise the Sun.\n",
			wantErr: "",
		},
		{
			name:    "test ad-hoc replacements",
			args:    []string{"-o", "MY_TEMPLATE {{ .at_message }}.", "-r", "@ at_", "-b", "4096"},
			input:   "{ \"@message\": \"Praise the Sun\" }",
			wantW:   "MY_TEMPLATE Praise the Sun.\n",
			wantErr: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() { viper.Reset() }()
			RootCmd.SetArgs(tt.args)

			tempDir := t.TempDir()
			initRealUserConfigDirFunc := configuration.UserConfigDir
			loadRealUserConfigDirFunc := load.UserConfigDir
			realExit := exitCmd

			mockUserConfigDirFunc := func() (string, error) {
				return tempDir, nil
			}
			configuration.UserConfigDir = mockUserConfigDirFunc
			load.UserConfigDir = mockUserConfigDirFunc
			exitCmd = func(code int) {
				fmt.Printf("exited with code %d\n", code)
			}

			defer func() {
				configuration.UserConfigDir = initRealUserConfigDirFunc
				load.UserConfigDir = loadRealUserConfigDirFunc
				exitCmd = realExit
			}()

			gotW, gotErr := executeCommandWithInput(t, tt.input)

			if gotW != tt.wantW {
				t.Errorf("executeCommandWithInput gotW = '%v', want '%v'", gotW, tt.wantW)
			}

			if tt.wantErr != "" && !strings.HasSuffix(gotErr, tt.wantErr) {
				t.Errorf("executeCommandWithInput gotErr = '%v', want it to end with '%v'", gotErr, tt.wantErr)
			}
		})
	}
}

func executeCommandWithInput(t *testing.T, input string) (string, string) {
	inputBytes := []byte(input)
	r := prepareStdIn(t, inputBytes)
	readStdOut, wStdout := prepareStdOut(t)
	readStdErr, wStderr := prepareStdOut(t)

	stdin := os.Stdin
	stdout := os.Stdout
	stderr := os.Stderr
	defer func() {
		os.Stdin = stdin
		os.Stdout = stdout
		os.Stderr = stderr
	}()
	os.Stdin = r
	os.Stdout = wStdout
	os.Stderr = wStderr

	err := RootCmd.Execute()
	if err != nil {
		t.Fatal(err)
	}
	err = wStdout.Close()
	if err != nil {
		t.Fatal(err)
	}
	err = wStderr.Close()
	if err != nil {
		t.Fatal(err)
	}

	output, err := ioutil.ReadAll(readStdOut)
	if err != nil {
		t.Fatal(err)
	}
	outputErr, err := ioutil.ReadAll(readStdErr)
	if err != nil {
		t.Fatal(err)
	}
	return string(output), string(outputErr)
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
