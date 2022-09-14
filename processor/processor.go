package processor

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/mdreem/json2line/common"
	"io"
	"sort"
	"strings"
	"text/template"
)

const InitialBufferSize = 1048576

func ProcessInput(reader io.Reader, writer io.Writer, t *template.Template, replacements map[string]string, bufferSize int) error {
	scanner := bufio.NewScanner(bufio.NewReader(reader))
	scanner.Buffer(make([]byte, 0, bufferSize), bufferSize)

	for scanner.Scan() {
		_, e := fmt.Fprintln(writer, processJSON(scanner.Text(), t, replacements))
		if e != nil {
			return e
		}
	}
	if scanner.Err() != nil {
		return scanner.Err()
	}
	return nil
}

func processJSON(input string, t *template.Template, replacements map[string]string) string {
	var parsedJSON map[string]interface{}
	err := json.Unmarshal([]byte(input), &parsedJSON)
	if err != nil {
		return input
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
		common.PrintInformationf("could not parse template line: %v", err)
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
			common.PrintInformationf("unknown type")
		}
	}
}
