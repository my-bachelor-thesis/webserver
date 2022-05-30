package user_solutions_with_tests

import (
	"webserver/internal/postgres"
	"webserver/internal/postgres/rdg/tests"
	"webserver/internal/postgres/rdg/user_solutions"
)

func GetByLanguage(language string, taskId int, userId int) (*UserSolutionsWithTests, error) {
	res := &UserSolutionsWithTests{}

	// TODO in transaction

	solutionsStatement := `
	select 
		us.id,
		to_char(us.last_modified, 'DD.MM.YY, HH24:MI:SS'),
		us.name,
		us.public,
		coalesce((select ust.test_id from user_solutions_tests ust where ust.user_solution_id = us.id and ust.user_id = $1), 0) as test_id
	from user_solutions us where us.language = $2 and us.task_id = $3 and (us.user_id = $4 or us.public) order by us.last_modified desc`

	rows, err := postgres.GetPool().Query(postgres.GetCtx(), solutionsStatement, userId, language, taskId, userId)
	if err != nil {
		return nil, err
	}
	var us Solutions
	var testId int
	for rows.Next() {
		if err = rows.Scan(&us.Id, &us.LastModified, &us.Name, &us.Public, &testId); err != nil {
			return nil, err
		}
		res.Solutions = append(res.Solutions, &Solutions{
			UserSolution: user_solutions.UserSolution{Id: us.Id, LastModified: us.LastModified, Name: us.Name, Public: us.Public},
			TestId:       testId,
		})
	}

	testsStatement := `
	select 
		id,
		to_char(last_modified, 'DD.MM.YY, HH24:MI:SS'),
		final,
		name,
		public
	from tests where language = $1 and task_id = $2 and (user_id = $3 or public = true) order by final desc, last_modified desc`

	rows, err = postgres.GetPool().Query(postgres.GetCtx(), testsStatement, language, taskId, userId)
	if err != nil {
		return nil, err
	}

	var test tests.Test
	for rows.Next() {
		if err = rows.Scan(&test.Id, &test.LastModified, &test.Final, &test.Name, &test.Public); err != nil {
			return nil, err
		}
		res.Tests = append(res.Tests, &tests.Test{
			Id: test.Id, LastModified: test.LastModified, Final: test.Final, Name: test.Name, Public: test.Public})
	}

	return res, err
}
