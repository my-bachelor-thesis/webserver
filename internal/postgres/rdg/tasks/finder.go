package tasks

import (
	"fmt"
	"github.com/jackc/pgx/v4"
	"webserver/internal/postgres"
	"webserver/pkg/postgresutil"
)

func GetById(id int) (*Task, error) {
	if id == 0 {
		return nil, postgresutil.ErrNoRowsInResult
	}
	statement := fmt.Sprintf("select %s from tasks where id = $1", allFieldReplacedTimestamp)
	task := Task{}
	err := load(postgres.GetPool().QueryRow(postgres.GetCtx(), statement, id), &task)
	return &task, err
}

func GetApprovedAndPublished() ([]*Task, error) {
	return getManyWithConditions("is_published = true and approver_id != 0")
}

func GetUnapproved() ([]*Task, error) {
	return getManyWithConditions("is_published = true and approver_id = 0")
}

func GetUnpublished() ([]*Task, error) {
	return getManyWithConditions("is_published = false")
}

func getManyWithConditions(conditions string) ([]*Task, error) {
	statement := fmt.Sprintf("select %s from tasks where %s", allFieldReplacedTimestamp, conditions)
	rows, err := postgres.GetPool().Query(postgres.GetCtx(), statement)
	if err != nil {
		return nil, err
	}
	var tasks []*Task
	for rows.Next() {
		task := Task{}
		if err = load(rows, &task); err != nil {
			return nil, err
		}
		if task.Id == 0 {
			continue
		}
		tasks = append(tasks, &task)
	}
	if len(tasks) == 0 {
		return nil, postgresutil.ErrNoRowsInResult
	}
	return tasks, err
}

func load(qr pgx.Row, task *Task) error {
	err := qr.Scan(&task.Id, &task.AuthorId, &task.ApproverId, &task.Title,
		&task.Difficulty, &task.IsPublished, &task.AddedOn, &task.Text)
	return err
}
