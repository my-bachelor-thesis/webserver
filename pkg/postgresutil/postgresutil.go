package postgresutil

import (
	"fmt"
	"strings"
)

func GeneratePlaceholder(fields string) (placeHolders string) {
	n := strings.Count(fields, ",") + 1
	for i := 1; i <= n; i++ {
		if i != n {
			placeHolders += fmt.Sprintf("$%d, ", i)
		} else {
			placeHolders += fmt.Sprintf("$%d", i)
		}
	}
	return
}

func CallToCharOnTimestamp(s, fieldToReplace string) string {
	return strings.Replace(s, fieldToReplace, fmt.Sprintf("to_char(%s, 'DD.MM.YY, HH24:MI:SS')", fieldToReplace), 1)
}

