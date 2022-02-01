package postgresutil

import (
	"fmt"
	"strings"
)

func GeneratePlaceholders(fields string) (placeHolders string) {
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

func GeneratePlaceholdersAndReplace(fields string, replace map[int]string) (placeHolders string) {
	n := strings.Count(fields, ",") + 1
	var numberOfReplaced int
	for i := 1; i <= n; i++ {
		if value, ok := replace[i-1]; ok {
			placeHolders += fmt.Sprintf("%s, ", value)
			numberOfReplaced++
			continue
		}
		placeHolders += fmt.Sprintf("$%d, ", i-numberOfReplaced)
	}
	return placeHolders[:len(placeHolders)-2]
}

func CallToCharOnTimestamp(s, fieldToReplace string) string {
	return strings.Replace(s, fieldToReplace, fmt.Sprintf("to_char(%s, 'DD.MM.YY, HH24:MI:SS')", fieldToReplace), 1)
}
