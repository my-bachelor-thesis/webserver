package transaction_scripts

import (
	"github.com/jackc/pgx/v4"
	"github.com/labstack/echo/v4"
	"webserver/internal/jwt"
	"webserver/internal/postgres"
	"webserver/internal/postgres/rdg/task_with_solutions_and_tests"
	"webserver/internal/postgres/rdg/tasks"
	"webserver/internal/postgres/rdg/tests"
	"webserver/internal/postgres/rdg/user_solutions"
	"webserver/internal/postgres/rdg/users"
)

func PublishTask(c echo.Context, claims *jwt.CustomClaims) error {
	conn, tx, err := getConnectionFromPoolAndStartTrans(pgx.RepeatableRead)
	if err != nil {
		return err
	}
	defer tx.Rollback(postgres.GetCtx())
	defer conn.Release()

	_, task, err := bindAndFindWithUserId(tx, c, tasks.GetByIdAndAuthorId, claims.UserId)
	if err != nil {
		return err
	}

	if claims.IsAdmin {
		task.ApproverId = claims.UserId
		if err := task.ApproveAndPublish(tx); err != nil {
			return err
		}
	} else if err := task.Publish(tx); err != nil {
		return err
	}

	return tx.Commit(postgres.GetCtx())
}

func UnpublishTask(c echo.Context, claims *jwt.CustomClaims) error {
	conn, tx, err := getConnectionFromPoolAndStartTrans(pgx.RepeatableRead)
	if err != nil {
		return err
	}
	defer tx.Rollback(postgres.GetCtx())
	defer conn.Release()

	_, task, err := bindAndFindWithUserId(tx, c, tasks.GetByIdAndAuthorId, claims.UserId)
	if err != nil {
		return err
	}

	if err := task.Unpublish(tx); err != nil {
		return err
	}

	return tx.Commit(postgres.GetCtx())
}

func DeleteTask(c echo.Context, claims *jwt.CustomClaims) error {
	conn, tx, err := getConnectionFromPoolAndStartTrans(pgx.RepeatableRead)
	if err != nil {
		return err
	}
	defer tx.Rollback(postgres.GetCtx())
	defer conn.Release()

	_, task, err := bindAndFindWithUserId(tx, c, tasks.GetByIdAndAuthorId, claims.UserId)
	if err != nil {
		return err
	}

	if err := task.Delete(tx); err != nil {
		return err
	}

	return tx.Commit(postgres.GetCtx())
}

func ApproveTask(c echo.Context, claims *jwt.CustomClaims) error {
	conn, tx, err := getConnectionFromPoolAndStartTrans(pgx.RepeatableRead)
	if err != nil {
		return err
	}
	defer tx.Rollback(postgres.GetCtx())
	defer conn.Release()

	_, task, err := bindAndFind(tx, c, tasks.GetById)
	if err != nil {
		return err
	}

	task.ApproverId = claims.UserId
	if err := task.Approve(tx); err != nil {
		return err
	}

	return tx.Commit(postgres.GetCtx())
}

func DenyTask(claims *jwt.CustomClaims, taskId, authorId int) (task *tasks.Task, user *users.User, admin *users.User, err error) {
	conn, tx, err := getConnectionFromPoolAndStartTrans(pgx.RepeatableRead)
	if err != nil {
		return nil, nil, nil, err
	}
	defer tx.Rollback(postgres.GetCtx())
	defer conn.Release()

	task, err = tasks.GetById(tx, taskId)
	if err != nil {
		return nil, nil, nil, err
	}

	admin, err = users.GetById(tx, claims.UserId)
	if err != nil {
		return nil, nil, nil, err
	}

	user, err = users.GetById(tx, authorId)
	if err != nil {
		return nil, nil, nil, err
	}

	if err := task.Unapprove(tx); err != nil {
		return nil, nil, nil, err
	}

	return task, user, admin, tx.Commit(postgres.GetCtx())
}

func AddTask(claims *jwt.CustomClaims, request *task_with_solutions_and_tests.TaskWithSolutionsAndTests) error {

	conn, tx, err := getConnectionFromPoolAndStartTrans(pgx.RepeatableRead)
	if err != nil {
		return err
	}
	defer tx.Rollback(postgres.GetCtx())
	defer conn.Release()

	// insert task
	task := tasks.Task{
		AuthorId:   claims.UserId,
		Title:      request.Title,
		Difficulty: request.Difficulty,
		Text:       request.Description,
	}

	// if updating
	if request.TaskId != 0 {
		task.Id = request.TaskId
		if err = task.UpdateTitleDifficultyAndText(tx); err != nil {
			return err
		}
		if err := deleteAllPublicOrFinal(tx, task.Id); err != nil {
			return err
		}
	} else {
		if err = task.Insert(tx); err != nil {
			return err
		}
	}

	if len(request.PublicSolutions) > 0 {
		var publicSolutions []*user_solutions.UserSolution
		for _, solution := range request.PublicSolutions {
			var u user_solutions.UserSolution
			u.UserId = claims.UserId
			u.TaskId = task.Id
			u.Language = solution.Language
			u.Name = solution.Name
			u.Public = true
			u.Code = solution.Code
			publicSolutions = append(publicSolutions, &u)
		}
		if err = user_solutions.InsertMany(tx, publicSolutions); err != nil {
			return err
		}
	}

	fillTest := func(newTest *tests.Test, testFromRequest *task_with_solutions_and_tests.NameAndCode) {
		newTest.Name = testFromRequest.Name
		newTest.Public = true
		newTest.UserId = claims.UserId
		newTest.TaskId = task.Id
		newTest.Language = testFromRequest.Language
		newTest.Code = testFromRequest.Code
	}

	if len(request.PublicTests) > 0 {
		var publicTests []*tests.Test
		for _, test := range request.PublicTests {
			var t tests.Test
			fillTest(&t, test)
			publicTests = append(publicTests, &t)
		}
		if err = tests.InsertMany(tx, publicTests); err != nil {
			return err
		}
	}

	var finalTests []*tests.Test
	for _, test := range request.FinalTests {
		var t tests.Test
		fillTest(&t, test)
		t.Final = true
		t.Name = "Final"
		finalTests = append(finalTests, &t)
	}

	if err := tests.InsertMany(tx, finalTests); err != nil {
		return err
	}

	return tx.Commit(postgres.GetCtx())
}

func deleteAllPublicOrFinal(tx pgx.Tx, taskId int) error {
	statementDeleteTests := "delete from tests where task_id = $1 and (final = true or public = true)"
	if _, err := tx.Exec(postgres.GetCtx(), statementDeleteTests, taskId); err != nil {
		return err
	}

	statementDeleteUserSolutionsResults := `
	delete from user_solutions_results
	where user_solution_id in (select id from user_solutions where task_id = $1)`
	if _, err := tx.Exec(postgres.GetCtx(), statementDeleteUserSolutionsResults, taskId); err != nil {
		return err
	}

	statementDeleteLastOpened := "delete from last_opened where task_id = $1"
	if _, err := tx.Exec(postgres.GetCtx(), statementDeleteLastOpened, taskId); err != nil {
		return err
	}

	statementDeleteSolutions := "delete from user_solutions where task_id = $1 and public = true"
	_, err := tx.Exec(postgres.GetCtx(), statementDeleteSolutions, taskId)
	return err
}

func bindAndFindWithUserId[T any](tx postgres.PoolInterface, c echo.Context,
	getByIdFunc func(postgres.PoolInterface, int, int) (T, error), userId int) (*requestWithIdAndName, T, error) {

	var req requestWithIdAndName
	if err := c.Bind(&req); err != nil {
		return nil, *new(T), err
	}
	obj, err := getByIdFunc(tx, req.Id, userId)
	return &req, obj, err
}

func bindAndFind[T any](tx postgres.PoolInterface, c echo.Context,
	getByIdFunc func(postgres.PoolInterface, int) (T, error)) (*requestWithIdAndName, T, error) {
	var req requestWithIdAndName
	if err := c.Bind(&req); err != nil {
		return nil, *new(T), err
	}
	obj, err := getByIdFunc(tx, req.Id)
	return &req, obj, err
}

type requestWithIdAndName struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}
