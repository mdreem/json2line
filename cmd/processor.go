package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"
)

func ProcessInput(r io.Reader, w io.Writer) error {
	scanner := bufio.NewScanner(bufio.NewReader(r))
	for scanner.Scan() {
		_, e := fmt.Fprintln(w, processJSON(scanner.Text()))
		if e != nil {
			return e
		}
	}
	return nil
}

func processJSON(input string) string {
	var parsedJSON map[string]interface{}
	err := json.Unmarshal([]byte(input), &parsedJSON)
	if err != nil {
		fmt.Printf("could no parse line: %v", err)
	}

	var resultStrings []string

	appendValues(parsedJSON, &resultStrings)
	return strings.Join(resultStrings, " ")
}

func appendValues(parsedJSON map[string]interface{}, resultStrings *[]string) {
	keys := make([]string, 0)
	for k := range parsedJSON {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		element := parsedJSON[key]

		switch e := element.(type) {
		case string:
			*resultStrings = append(*resultStrings, e)
		case map[string]interface{}:
			appendValues(e, resultStrings)
		default:
			fmt.Println("unknown type")
		}
	}
}
