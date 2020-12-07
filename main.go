package main

import (
	"fmt"
	"github.com/mdreem/json2line/cmd"
	"os"
)

func main() {
	if isInputFromPipe() {
		err := cmd.ProcessInput(os.Stdin, os.Stdout)
		if err != nil {
			fmt.Printf("could not parse line: %v", err)
		}
	}
}

func isInputFromPipe() bool {
	stat, err := os.Stdin.Stat()
	if err != nil {
		panic(fmt.Errorf("could not read stat"))
	}
	return stat.Mode()&os.ModeCharDevice == 0
}
