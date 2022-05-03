package user_solutions_with_tests

import (
	"webserver/internal/postgres"
)

func GetByLanguage(language string, taskId int, userId int) (*UserSolutionsWithTests, error) {
	res := NewUserSolutionsWithTests()

	solutionsStatement := `
	select 
		id,
		to_char(last_modified, 'DD.MM.YY, HH24:MI:SS'),
		exit_code,
		name
	from user_solutions where user_id = $1 and language = $2 and task_id = $3 order by id`

	rows, err := postgres.GetPool().Query(postgres.GetCtx(), solutionsStatement, userId, language, taskId)
	if err != nil {
		return nil, err
	}
	var id int
	var lastModified string
	var exitCode int
	var name string
	for rows.Next() {
		if err = rows.Scan(&id, &lastModified, &exitCode, &name); err != nil {
			return nil, err
		}
		res.Solutions[id] = Solution{LastModified: lastModified, ExitCode: exitCode, Name: name}
	}

	testsStatement := `
	select 
		id,
		to_char(last_modified, 'DD.MM.YY, HH24:MI:SS'),
		final,
		name
	from tests where (user_id = $1 or final = true) and language = $2 and task_id = $3 order by id`

	rows, err = postgres.GetPool().Query(postgres.GetCtx(), testsStatement, userId, language, taskId)
	if err != nil {
		return nil, err
	}
	var final bool
	for rows.Next() {
		if err = rows.Scan(&id, &lastModified, &final, &name); err != nil {
			return nil, err
		}
		res.Tests[id] = Test{LastModified: lastModified, Final: final, Name: name}
	}

	return res, err
}
