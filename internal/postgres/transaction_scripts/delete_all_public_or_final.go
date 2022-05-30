package transaction_scripts

import (
	"webserver/internal/postgres"
)

func DeleteAllPublicOrFinal(taskId int) error {
	statementDeleteTests := "delete from tests where task_id = $1 and (final = true or public = true)"
	if _, err := postgres.GetPool().Exec(postgres.GetCtx(), statementDeleteTests, taskId); err != nil {
		return err
	}

	statementDeleteUserSolutionsResults := `
	delete from user_solutions_results
	where user_solution_id in (select id from user_solutions where task_id = $1)`
	if _, err := postgres.GetPool().Exec(postgres.GetCtx(), statementDeleteUserSolutionsResults, taskId); err != nil {
		return err
	}

	statementDeleteLastOpened := "delete from last_opened where task_id = $1"
	if _, err := postgres.GetPool().Exec(postgres.GetCtx(), statementDeleteLastOpened, taskId); err != nil {
		return err
	}

	statementDeleteSolutions := "delete from user_solutions where task_id = $1 and public = true"
	_, err := postgres.GetPool().Exec(postgres.GetCtx(), statementDeleteSolutions, taskId)
	return err
}
