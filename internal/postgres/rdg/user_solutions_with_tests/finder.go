package user_solutions_with_tests

import (
	"webserver/internal/postgres"
	"webserver/internal/postgres/rdg/tests"
	"webserver/internal/postgres/rdg/user_solutions"
)

func GetByLanguage(language string, taskId int, userId int) (*UserSolutionsWithTests, error) {
	res := NewUserSolutionsWithTests()

	solutionsStatement := `
	select 
		id,
		to_char(last_modified, 'DD.MM.YY, HH24:MI:SS'),
		exit_code,
		name,
		coalesce((select test_id from user_solutions_test_ids usti where usti.user_solution_id = us.id and usti.user_id = $1), 0) as test_id
	from user_solutions us where us.user_id = $2 and us.language = $3 and us.task_id = $4 order by us.last_modified desc`

	rows, err := postgres.GetPool().Query(postgres.GetCtx(), solutionsStatement, userId, userId, language, taskId)
	if err != nil {
		return nil, err
	}
	var us Solutions
	var testId int
	for rows.Next() {
		if err = rows.Scan(&us.Id, &us.LastModified, &us.ExitCode, &us.Name, &testId); err != nil {
			return nil, err
		}
		res.Solutions[us.Id] = Solutions{
			UserSolution: user_solutions.UserSolution{LastModified: us.LastModified, ExitCode: us.ExitCode, Name: us.Name},
			TestId:       testId,
		}
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

	var test tests.Test
	for rows.Next() {
		if err = rows.Scan(&test.Id, &test.LastModified, &test.Final, &test.Name); err != nil {
			return nil, err
		}
		res.Tests[test.Id] = tests.Test{LastModified: test.LastModified, Final: test.Final, Name: test.Name}
	}

	return res, err
}
