package cmd

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"text/template"
)

func ProcessInput(r io.Reader, w io.Writer, t *template.Template, replacements map[string]string) error {
	scanner := bufio.NewScanner(bufio.NewReader(r))
	for scanner.Scan() {
		_, e := fmt.Fprintln(w, processJSON(scanner.Text(), t, replacements))
		if e != nil {
			return e
		}
	}
	return nil
}

func processJSON(input string, t *template.Template, replacements map[string]string) string {
	var parsedJSON map[string]interface{}
	err := json.Unmarshal([]byte(input), &parsedJSON)
	if err != nil {
		printInformationf("could no parse line: %v", err)
	}
	replaceKeys(&parsedJSON, replacements)

	if t == nil {
		var resultStrings []string

		appendValues(parsedJSON, &resultStrings)
		return strings.Join(resultStrings, " ")
	}
	var buffer bytes.Buffer
	err = t.Execute(&buffer, parsedJSON)
	if err != nil {
		printInformationf("could no template line: %v", err)
	}
	return buffer.String()
}

func replaceKeys(data *map[string]interface{}, replacements map[string]string) {
	for k, v := range *data {
		delete(*data, k)

		mapValue, ok := v.(map[string]interface{})
		if ok {
			replaceKeys(&mapValue, replacements)
		}

		var keyWithReplacedCharacters = k
		for term, replacement := range replacements {
			keyWithReplacedCharacters = strings.ReplaceAll(keyWithReplacedCharacters, term, replacement)
		}
		(*data)[keyWithReplacedCharacters] = v
	}
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
			printInformationf("unknown type")
		}
	}
}

func printInformationf(format string, a ...interface{}) {
	_, err := fmt.Fprintf(os.Stderr, format, a...)
	if err != nil {
		panic(fmt.Errorf("could not print to stderr: %v", err))
	}
}
