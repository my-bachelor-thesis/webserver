package transaction_scripts

import (
	"github.com/jackc/pgx/v4"
	"webserver/internal/postgres"
	"webserver/internal/postgres/rdg/task_with_solutions_and_tests"
)

func GetTaskWithSolutionsAndTasksByTaskId(taskId, authorId int) (*task_with_solutions_and_tests.TaskWithSolutionsAndTests, error) {
	conn, tx, err := getConnectionFromPoolAndStartTrans(pgx.RepeatableRead)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(postgres.GetCtx())
	defer conn.Release()

	taskStatement := `
	select
		t.title,
		t.difficulty,
		t.text,
		t.is_published
	from tasks t where id = $1 and author_id = $2`
	task := task_with_solutions_and_tests.TaskWithSolutionsAndTests{}
	if err := tx.QueryRow(postgres.GetCtx(), taskStatement, taskId, authorId).Scan(&task.Title, &task.Difficulty,
		&task.Description, &task.IsPublished); err != nil {
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

	rows, err := tx.Query(postgres.GetCtx(), testsStatement, taskId)
	if err != nil {
		return nil, err
	}

	var final bool
	for rows.Next() {
		test := &task_with_solutions_and_tests.NameAndCode{}
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

	rows, err = tx.Query(postgres.GetCtx(), solutionStatement, taskId)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		solution := &task_with_solutions_and_tests.NameAndCode{}
		if err = rows.Scan(&solution.Name, &solution.Code, &solution.Language); err != nil {
			return nil, err
		}
		task.PublicSolutions = append(task.PublicSolutions, solution)
	}

	err = tx.Commit(postgres.GetCtx())
	return &task, err
}
