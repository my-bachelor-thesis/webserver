package transaction_scripts

import (
	"webserver/internal/postgres"
	"webserver/internal/postgres/rdg/tests"
	"webserver/internal/postgres/rdg/user_solutions"
	"webserver/internal/postgres/rdg/user_solutions_results"
)

type InsertedSolution struct {
	Id           int    `json:"id"`            // id of a solution that was inserted into the db
	LastModified string `json:"last_modified"` // last_modified of a solution that was inserted into the db
}

type InsertedTest struct {
	Id           int    `json:"id"`            // id of a test that was inserted into the db
	LastModified string `json:"last_modified"` // last_modified of a test that was inserted into the db
}

func EditorSaveTestAndSolution(test *tests.Test, solution *user_solutions.UserSolution,
	result *user_solutions_results.UserSolutionResult) (*InsertedSolution, *InsertedTest, error) {

	conn, tx, err := getConnectionFromPoolAndStartRegularTrans()
	if err != nil {
		return nil, nil, err
	}
	defer tx.Rollback(postgres.GetCtx())
	defer conn.Release()

	if err := test.Insert(tx); err != nil {
		return nil, nil, err
	}

	if err := solution.Insert(tx); err != nil {
		return nil, nil, err
	}

	result.TestId = test.Id
	result.UserSolutionId = solution.UserId

	if err := result.Insert(tx); err != nil {
		return nil, nil, err
	}

	err = tx.Commit(postgres.GetCtx())
	return &InsertedSolution{Id: solution.Id, LastModified: solution.LastModified}, &InsertedTest{Id: test.Id,
		LastModified: test.LastModified}, err
}

func EditorSaveTest(test *tests.Test, result *user_solutions_results.UserSolutionResult) (*InsertedTest, error) {
	conn, tx, err := getConnectionFromPoolAndStartRegularTrans()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(postgres.GetCtx())
	defer conn.Release()

	if err := test.Insert(tx); err != nil {
		return nil, err
	}

	result.TestId = test.Id
	if err := result.Insert(tx); err != nil {
		return nil, err
	}

	err = tx.Commit(postgres.GetCtx())
	return &InsertedTest{Id: test.Id, LastModified: test.LastModified}, err
}

func EditorSaveSolution(solution *user_solutions.UserSolution, result *user_solutions_results.UserSolutionResult) (*InsertedSolution, error) {
	conn, tx, err := getConnectionFromPoolAndStartRegularTrans()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(postgres.GetCtx())
	defer conn.Release()

	if err := solution.Insert(tx); err != nil {
		return nil, err
	}

	result.UserSolutionId = solution.Id
	if err := result.Insert(tx); err != nil {
		return nil, err
	}

	err = tx.Commit(postgres.GetCtx())
	return &InsertedSolution{Id: solution.Id, LastModified: solution.LastModified}, err
}
