package postgresutil

import (
	"errors"
	"fmt"
	"github.com/jackc/pgconn"
	"strings"
)

var ErrNoRowsInResult = errors.New("no rows in result set")

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
	return GeneratePlaceholdersAndReplaceFromIndex(fields, replace, 1)
}

func GeneratePlaceholdersAndReplaceFromIndex(fields string, replace map[int]string, index int) (placeHolders string) {
	n := strings.Count(fields, ",") + 1
	var numberOfReplaced int
	for i := index; i < n+index; i++ {
		if value, ok := replace[i-index]; ok {
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

func IsNoRowsInResultErr(err error) bool {
	if err == nil {
		return false
	}
	return err.Error() == ErrNoRowsInResult.Error()
}

func IsUniqueConstraintErr(err error) bool {
	if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
		return true
	}
	return false
}