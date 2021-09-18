package cmd

import (
	"github.com/mdreem/json2line/common"
	"github.com/spf13/viper"
	"os"
	"strings"
	"text/template"
)

func loadTemplate(formatter string) *template.Template {
	templates := viper.GetStringMapString("templates")
	formatString := templates[formatter]
	if formatString != "" {
		return createTemplate(formatter, formatString)
	}
	common.PrintInformationf("no formatter with the name '%s' defined\n", formatter)
	return nil
}

func loadReplacements(adhocReplacements []string) map[string]string {
	configuredReplacements := viper.GetStringMapString("replacements")

	var replacements = make(map[string]string)
	for k, v := range configuredReplacements {
		replacements[k] = v
	}

	for _, r := range adhocReplacements {
		replacement := strings.Fields(r)
		if len(replacement) != 2 {
			common.PrintInformationf("replacement is not of correct format '<FROM> <TO>' where both values are separated via whitespace")
		}

		replacements[replacement[0]] = replacement[1]
	}
	return replacements
}

func createTemplate(formatter string, formatString string) *template.Template {
	parse, err := template.New(formatter).Parse(formatString)
	if err != nil {
		common.PrintInformationf("could not parse template: %\n", err)
		os.Exit(1)
	}
	return parse
}
