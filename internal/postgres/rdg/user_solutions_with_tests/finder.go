package user_solutions_with_tests

import (
	"fmt"
	"webserver/internal/postgres"
)

func GetByLanguageAndTaskId(language string, taskId int) (*UserSolutionsWithTests, error) {
	var res UserSolutionsWithTests

	solutionsStatement := `
	select 
		id,
		to_char(last_modified, 'DD.MM.YY, HH24:MI:SS'),
		exit_code
	from user_solutions where language = $1 and task_id = $2 order by id`

	rows, err := postgres.GetPool().Query(postgres.GetCtx(), solutionsStatement, fmt.Sprintf("'%s'", language), taskId)
	if err != nil {
		return nil, err
	}
	var id int
	var date string
	var exitCode int
	for rows.Next() {
		if err = rows.Scan(&id, &date, &exitCode); err != nil {
			return nil, err
		}
		res.Solutions[id] = Solution{Date: date, ExitCode: exitCode}
	}

	testsStatement := `
	select 
		id,
		to_char(last_modified, 'DD.MM.YY, HH24:MI:SS')
	from tests where language = $1 and task_id = $2 order by id`

	rows, err = postgres.GetPool().Query(postgres.GetCtx(), testsStatement, fmt.Sprintf("'%s'", language), taskId)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		if err = rows.Scan(&id, &date); err != nil {
			return nil, err
		}
		res.Tests[id] = Test{Date: date}
	}

	return &res, err
}
