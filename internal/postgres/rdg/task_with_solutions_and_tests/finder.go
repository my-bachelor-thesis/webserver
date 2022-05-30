package task_with_solutions_and_tests

import "webserver/internal/postgres"

func GetByTaskId(taskId, authorId int) (*TaskWithSolutionsAndTests, error) {
	// TODO: in transaction

	taskStatement := `
	select
		t.title,
		t.difficulty,
		t.text
	from tasks t where id = $1 and author_id = $2`
	task := TaskWithSolutionsAndTests{}
	if err := postgres.GetPool().QueryRow(postgres.GetCtx(), taskStatement, taskId, authorId).Scan(&task.Title, &task.Difficulty,
		&task.Description); err != nil {
		return nil, err
	}

	testsStatement := `
	select
		t.name,
		t.code,
		t.language,
		t.final
	from tests t
	where t.task_id = $1 and (t.final = true or t.public = true)`

	rows, err := postgres.GetPool().Query(postgres.GetCtx(), testsStatement, taskId)
	if err != nil {
		return nil, err
	}

	var final bool
	for rows.Next() {
		test := &NameAndCode{}
		if err = rows.Scan(&test.Name, &test.Code, &test.Language, &final); err != nil {
			return nil, err
		}
		if final {
			task.FinalTests = append(task.FinalTests, test)
		} else {
			task.PublicTests = append(task.PublicTests, test)
		}
	}

	solutionStatement := `
	select
		us.name,
		us.code,
		us.language
	from user_solutions us
	where us.task_id = $1 and us.public = true`

	rows, err = postgres.GetPool().Query(postgres.GetCtx(), solutionStatement, taskId)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		solution := &NameAndCode{}
		if err = rows.Scan(&solution.Name, &solution.Code, &solution.Language); err != nil {
			return nil, err
		}
		task.PublicSolutions = append(task.PublicSolutions, solution)
	}

	return &task, nil
}
